package probe

import (
	"context"

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
	}
)

// NewProbe initializes and returns a new Probe.
func NewProbe(ctx context.Context, work chan Runner, log *zerolog.Logger) *Probe {
	p := &Probe{
		log:  log,
		ctx:  ctx,
		work: work,
	}
	p.Work()
	return p
}

// Working returns the status of the Probe.
func (p *Probe) Working() bool {
	return p.working
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
		for {
			select {
			case <-p.childCtx.Done():
				// the context is done, exit
				p.working = false
				close(p.done)
				return
			case runner := <-p.work:
				runner()
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
