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
		log     *zerolog.Logger
		ctx     context.Context
		cancel  context.CancelFunc
		work    <-chan Runner
		done    chan struct{}
		working bool
	}
)

// NewProbe initializes and returns a new Probe.
func NewProbe(ctx context.Context, work <-chan Runner, log *zerolog.Logger) *Probe {
	childCtx, cancel := context.WithCancel(ctx)
	return &Probe{
		log:    log,
		ctx:    childCtx,
		cancel: cancel,
		work:   work,
		done:   make(chan struct{}),
	}
}

// Working returns the status of the Probe.
func (p *Probe) Working() bool {
	return p.working
}

// Work is the main event loop for the Probe. Work will start a new goroutine.
func (p *Probe) Work() {
	p.working = true
	go func() {
		for {
			select {
			case <-p.ctx.Done():
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

// Stop will stop the Probe from doing further work. Stop blocks if wait is true until all work is complete.
func (p *Probe) Stop(wait bool) {
	p.cancel()
	if wait {
		<-p.done
	}
}
