package pool

import (
	"context"
	"log/slog"
	"testing"

	"github.com/amplify-security/probe/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	// MockLogHandler for testing getLogHandler.
	MockLogHandler struct {
		mock.Mock
	}
)

// Enabled implementation of slog.Handler for pool.MockLogHandler.
// Returns false for all log levels.
func (h *MockLogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle implementation of slog.Handler for pool.MockLogHandler.
func (h *MockLogHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs implementation of slog.Handler for pool.MockLogHandler.
func (h *MockLogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup implementation of slog.Handler for pool.MockLogHandler.
func (h *MockLogHandler) WithGroup(_ string) slog.Handler {
	return h
}

func TestPoolConfig_getLogHandler(t *testing.T) {
	cases := []struct {
		h   slog.Handler
		msg string
	}{
		{
			h:   &MockLogHandler{},
			msg: "getLogHandler -> MockLogHandler",
		},
		{
			h: nil,
			msg: "getLogHandler -> NoopLogHandler",
		},
	}
	for _, c := range cases {
		cfg := &PoolConfig{
			LogHandler: c.h,
		}
		if c.h != nil {
			assert.Equal(t, c.h, cfg.getLogHandler(), c.msg)
		} else {
			assert.IsType(t, &logging.NoopLogHandler{}, cfg.getLogHandler(), c.msg)
		}
	}
}

func TestPoolConfig_getCtx(t *testing.T) {
	cases := []struct {
		ctx context.Context
		msg string
	}{
		{
			ctx: context.TODO(),
			msg: "getCtx -> context.TODO",
		},
		{
			ctx: nil,
			msg: "getCtx -> context.Background",
		},
	}
	for _, c := range cases {
		cfg := &PoolConfig{
			Ctx: c.ctx,
		}
		if c.ctx != nil {
			assert.Equal(t, c.ctx, cfg.getCtx(), c.msg)
		} else {
			assert.Equal(t, context.Background(), cfg.getCtx(), c.msg)
		}
	}
}

func TestPoolConfig_getSize(t *testing.T) {
	cases := []struct {
		size int
		msg  string
	}{
		{
			size: 1,
			msg:  "getSize -> 1",
		},
		{
			size: 0,
			msg:  "getSize -> DefaultPoolSize",
		},
	}
	for _, c := range cases {
		cfg := &PoolConfig{
			Size: c.size,
		}
		if c.size != 0 {
			assert.Equal(t, c.size, cfg.getSize(), c.msg)
		} else {
			assert.Equal(t, DefaultPoolSize, cfg.getSize(), c.msg)
		}
	}
}

func TestPoolConfig_getBufferSize(t *testing.T) {
	cases := []struct {
		size int
		msg  string
	}{
		{
			size: 1,
			msg:  "getBufferSize -> 1",
		},
		{
			size: 0,
			msg:  "getBufferSize -> DefaultBufferSize",
		},
	}
	for _, c := range cases {
		cfg := &PoolConfig{
			BufferSize: c.size,
		}
		if c.size != 0 {
			assert.Equal(t, c.size, cfg.getBufferSize(), c.msg)
		} else {
			assert.Equal(t, DefaultBufferSize, cfg.getBufferSize(), c.msg)
		}
	}
}
