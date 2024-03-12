package probe

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/amplify-security/probe/logging"
)

type (
	// ProbeConfig is a struct for passing configuration data to a new Probe.
	ProbeConfig struct {
		LogHandler slog.Handler    // Handler to use for probe logging. If empty, probe.NoopHandler will be used.
		Ctx        context.Context // Context to use for the probe. If empty, context.Background will be used.
		WorkChan   chan Runner     // Channel to use for work. If empty, a new channel will be created.
		RunningCtr *atomic.Int32   // Running counter to increment when this probe is running.
		IdleCtr    *atomic.Int32   // Idle counter to increment when this probe is idle.
		WaitGroup  *sync.WaitGroup // WaitGroup to use for the probe.
	}
)

// getLogHandler returns the log handler to use for the Probe.
func (c *ProbeConfig) getLogHandler() slog.Handler {
	if c.LogHandler == nil {
		return &logging.NoopLogHandler{}
	}
	return c.LogHandler
}

// getCtx returns the context to use for the Probe.
func (c *ProbeConfig) getCtx() context.Context {
	if c.Ctx == nil {
		return context.Background()
	}
	return c.Ctx
}

// getWorkChan returns the channel to use for the Probe.
func (c *ProbeConfig) getWorkChan() chan Runner {
	if c.WorkChan == nil {
		return make(chan Runner)
	}
	return c.WorkChan
}

// getRunningCtr returns the running counter to use for the Probe.
func (c *ProbeConfig) getRunningCtr() *atomic.Int32 {
	if c.RunningCtr == nil {
		return new(atomic.Int32)
	}
	return c.RunningCtr
}

// getIdleCtr returns the idle counter to use for the Probe.
func (c *ProbeConfig) getIdleCtr() *atomic.Int32 {
	if c.IdleCtr == nil {
		return new(atomic.Int32)
	}
	return c.IdleCtr
}

// getWaitGroup returns the WaitGroup to use for the Probe.
func (c *ProbeConfig) getWaitGroup() *sync.WaitGroup {
	if c.WaitGroup == nil {
		return new(sync.WaitGroup)
	}
	return c.WaitGroup
}
