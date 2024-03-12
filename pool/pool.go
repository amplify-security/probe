package pool

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/amplify-security/probe"
)

type (
	// Pool is a congigurable collection of Probes that run functions on available goroutines.
	Pool struct {
		logHandler slog.Handler
		log        *slog.Logger
		ctx        context.Context
		cancel     context.CancelFunc
		work       chan probe.Runner
		started    bool
		runningCtr *atomic.Int32
		idleCtr    *atomic.Int32
		waitGroup  *sync.WaitGroup
		size       int
		probes     []*probe.Probe
	}
)

// NewPool initializes and returns a new Pool.
func NewPool(cfg *PoolConfig) *Pool {
	logHandler := cfg.getLogHandler()
	log := slog.New(cfg.getLogHandler()).With("source", "probe.Pool")
	childCtx, cancel := context.WithCancel(cfg.getCtx())
	work := make(chan probe.Runner, cfg.getBufferSize())
	p := &Pool{
		logHandler: logHandler,
		log:        log,
		ctx:        childCtx,
		cancel:     cancel,
		work:       work,
		runningCtr: new(atomic.Int32),
		idleCtr:    new(atomic.Int32),
		waitGroup:  new(sync.WaitGroup),
		size:       cfg.getSize(),
		probes:     make([]*probe.Probe, 0, cfg.getSize()),
	}
	p.Start()
	return p
}

// Start starts the Pool.
func (p *Pool) Start() {
	if p.started {
		p.log.Info("received start request, but pool is already started")
		return
	}
	p.log.Info("starting pool")
	for _, p := range p.probes {
		// run all existing probes on restarts
		p.Run()
	}
	if len(p.probes) == 0 {
		// create all probes for new pools
		for range p.size {
			p.probes = append(p.probes, probe.NewProbe(&probe.ProbeConfig{
				LogHandler: p.logHandler,
				Ctx:        p.ctx,
				WorkChan:   p.work,
				RunningCtr: p.runningCtr,
				IdleCtr:    p.idleCtr,
				WaitGroup:  p.waitGroup,
			}))
		}
	}
	p.started = true
}

// Stop stops the Pool.
func (p *Pool) Stop(wait bool) {
	if !p.started {
		p.log.Info("received stop request, but pool is not started")
		return
	}
	p.log.Info("stopping pool")
	p.cancel()
	if wait {
		p.waitGroup.Wait()
	}
	p.started = false
}

// Run executes a probe.Runner on a Probe in the Pool.
func (p *Pool) Run(r probe.Runner) {
	p.work <- r
}

// Idle returns the number of idle Probes in the Pool.
func (p *Pool) Idle() int {
	return int(p.idleCtr.Load())
}

// Running returns the number of running Probes in the Pool.
func (p *Pool) Running() int {
	return int(p.runningCtr.Load())
}
