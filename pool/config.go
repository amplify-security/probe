package pool

import (
	"context"
	"log/slog"

	"github.com/amplify-security/probe/logging"
)

const (
	DefaultPoolSize   = 8  // DefaultPoolSize is the default size of the pool.
	DefaultBufferSize = 64 // DefaultBufferSize is the default size of the work channel buffer.
)

type (
	// PoolConfig is a struct for passing configuration data to a new Pool.
	PoolConfig struct {
		LogHandler slog.Handler    // Handler to use for pool logging. If empty, probe.NoopHandler will be used.
		Ctx        context.Context // Context to use for the pool. If empty, context.Background will be used.
		Size       int             // Size of the pool. Default pool size is 8.
		BufferSize int             // Size of the work channel buffer. Default buffer size is 64.
	}
)

// getLogHandler returns the log handler to use for the Pool.
func (c *PoolConfig) getLogHandler() slog.Handler {
	if c.LogHandler == nil {
		return &logging.NoopLogHandler{}
	}
	return c.LogHandler
}

// getCtx returns the context to use for the Pool.
func (c *PoolConfig) getCtx() context.Context {
	if c.Ctx == nil {
		return context.Background()
	}
	return c.Ctx
}

// getSize returns the size to use for the Pool.
func (c *PoolConfig) getSize() int {
	if c.Size == 0 {
		return DefaultPoolSize
	}
	return c.Size
}

// getBufferSize returns the buffer size to use for the Pool.
func (c *PoolConfig) getBufferSize() int {
	if c.BufferSize == 0 {
		return DefaultBufferSize
	}
	return c.BufferSize
}
