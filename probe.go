package probe

import (
	"context"
	"encoding/hex"
	"io"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// Runner function type.
	Runner func()

	// Probe is a helper that runs functions on a separate goroutine.
	Probe struct {
		log        *slog.Logger
		ctx        context.Context
		childCtx   context.Context
		cancel     context.CancelFunc
		work       chan Runner
		done       chan struct{}
		running    *atomic.Bool
		runningCtr *atomic.Int32
		idle       *atomic.Bool
		idleCtr    *atomic.Int32
		waitGroup  *sync.WaitGroup
		id         string
	}
)

var (
	randStream io.Reader = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// getID returns a random hex string to be used as a unique identifier for a Probe.
func getID(log *slog.Logger) string {
	buf := make([]byte, 6)
	if _, err := randStream.Read(buf); err != nil {
		log.Error("failed to generate probe ID", "error", err)
		panic(err)
	}
	return hex.EncodeToString(buf)[:6]
}

// NewProbe initializes and returns a new Probe.
func NewProbe(cfg *ProbeConfig) *Probe {
	log := slog.New(cfg.getLogHandler())
	id := getID(log)
	ctxLogger := log.With("id", id, "source", "probe.Probe")
	running := new(atomic.Bool)
	running.Store(false)
	idle := new(atomic.Bool)
	idle.Store(true)
	p := &Probe{
		log:        ctxLogger,
		ctx:        cfg.getCtx(),
		work:       cfg.getWorkChan(),
		running:    running,
		runningCtr: cfg.getRunningCtr(),
		idle:       idle,
		idleCtr:    cfg.getIdleCtr(),
		waitGroup:  cfg.getWaitGroup(),
		id:         id,
	}
	p.Run()
	return p
}

// ID returns the unique identifier of the Probe.
func (p *Probe) ID() string {
	return p.id
}

// Running returns the status of the Probe: true if work event loop is running.
func (p *Probe) Running() bool {
	return p.running.Load()
}

// Idle returns the status of the Probe: true if the Probe is Working but has no current work to execute.
func (p *Probe) Idle() bool {
	return p.idle.Load()
}

// WorkChan returns the channel used for work events.
func (p *Probe) WorkChan() chan Runner {
	return p.work
}

// Run is the main event loop for the Probe. Run will start a new goroutine.
func (p *Probe) Run() {
	if p.Running() {
		return
	}
	// create a new cancelable child context only to be used by this goroutine
	p.childCtx, p.cancel = context.WithCancel(p.ctx)
	p.waitGroup.Add(1)
	p.done = make(chan struct{})
	go func() {
		p.log.Debug("starting event loop")
		defer p.waitGroup.Done()
		p.running.Store(true)
		p.runningCtr.Add(1)
		p.idleCtr.Add(1)
		for {
			select {
			case <-p.childCtx.Done():
				// the context is done, exit
				p.log.Debug("shutting down")
				p.running.Store(false)
				p.idle.Store(true)
				p.runningCtr.Add(-1)
				close(p.done)
				return
			case runner := <-p.work:
				p.idle.Store(false)
				p.idleCtr.Add(-1)
				runner()
				p.idle.Store(true)
				p.idleCtr.Add(1)
			}
		}
	}()
}

// Stop will stop the Probe from doing further work. Stop blocks if wait is true until current work is complete.
func (p *Probe) Stop(wait bool) {
	if !p.Running() {
		return
	}
	p.cancel()
	if wait {
		<-p.done
		return
	}
}
