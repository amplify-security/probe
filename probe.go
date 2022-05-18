package probe

import (
	"context"
	"encoding/hex"
	"io"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
)

type (
	// Runner function type.
	Runner func()

	// Probe is a helper that runs functions on a separate goroutine.
	Probe struct {
		log      *zerolog.Logger
		ctx      context.Context
		childCtx context.Context
		cancel   context.CancelFunc
		work     chan Runner
		done     chan struct{}
		working  bool
		idle     bool
		id       string
	}
)

var (
	randStream io.Reader = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// getID returns a random hex string to be used as a unique identifier for a Probe.
func getID(log *zerolog.Logger) string {
	buf := make([]byte, 6)
	if _, err := randStream.Read(buf); err != nil {
		log.Panic().Err(err).Msg("unable to generate probe id")
	}
	return hex.EncodeToString(buf)[:6]
}

// NewProbe initializes and returns a new Probe.
func NewProbe(ctx context.Context, work chan Runner, log *zerolog.Logger) *Probe {
	id := getID(log)
	ctxLogger := log.With().Str("id", id).Str("source", "probe.Probe").Logger()
	p := &Probe{
		log:  &ctxLogger,
		ctx:  ctx,
		work: work,
		idle: true,
		id:   id,
	}
	p.Work()
	return p
}

// ID returns the unique identifier of the Probe.
func (p *Probe) ID() string {
	return p.id
}

// Working returns the status of the Probe: true if work event loop is running.
func (p *Probe) Working() bool {
	return p.working
}

// Idle returns the status of the Probe: true if the Probe is Working but has no current work to execute.
func (p *Probe) Idle() bool {
	return p.idle
}

// Work is the main event loop for the Probe. Work will start a new goroutine.
func (p *Probe) Work() {
	if p.Working() {
		return
	}
	// create a new cancelable child context only to be used by this goroutine
	p.childCtx, p.cancel = context.WithCancel(p.ctx)
	p.working = true
	p.done = make(chan struct{})
	go func() {
		p.log.Debug().Msg("starting work event loop")
		for {
			select {
			case <-p.childCtx.Done():
				// the context is done, exit
				p.log.Debug().Msg("shutting down")
				p.working = false
				close(p.done)
				return
			case runner := <-p.work:
				p.idle = false
				runner()
				p.idle = true
			}
		}
	}()
}

// Stop will stop the Probe from doing further work. Stop blocks if wait is true until current work is complete.
func (p *Probe) Stop(wait bool) {
	if !p.Working() {
		return
	}
	p.cancel()
	if wait {
		<-p.done
		return
	}
}
