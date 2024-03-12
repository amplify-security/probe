package probe

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
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

func TestProbeConfig_getLogHandler(t *testing.T) {
	cases := []struct {
		h   slog.Handler
		msg string
	}{
		{
			h:   &MockLogHandler{},
			msg: "getLogHandler -> MockLogHandler",
		},
		{
			h:   nil,
			msg: "getLogHandler -> NoopLogHandler",
		},
	}
	for _, c := range cases {
		cfg := &ProbeConfig{
			LogHandler: c.h,
		}
		if c.h != nil {
			assert.Equal(t, c.h, cfg.getLogHandler(), c.msg)
		} else {
			assert.IsType(t, &logging.NoopLogHandler{}, cfg.getLogHandler(), c.msg)
		}
	}
}

func TestProbeConfig_getCtx(t *testing.T) {
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
		cfg := &ProbeConfig{
			Ctx: c.ctx,
		}
		if c.ctx != nil {
			assert.Equal(t, c.ctx, cfg.getCtx(), c.msg)
		} else {
			assert.Equal(t, context.Background(), cfg.getCtx(), c.msg)
		}
	}
}

func TestProbeConfig_getWorkChan(t *testing.T) {
	cases := []struct {
		ch  chan Runner
		msg string
	}{
		{
			ch:  make(chan Runner),
			msg: "getWorkChan(ch) -> make(chan Runner)",
		},
		{
			ch:  nil,
			msg: "getWorkChan(nil) -> make(chan Runner)",
		},
	}
	for _, c := range cases {
		cfg := &ProbeConfig{
			WorkChan: c.ch,
		}
		if c.ch != nil {
			assert.Equal(t, c.ch, cfg.getWorkChan(), c.msg)
		} else {
			assert.NotNil(t, cfg.getWorkChan(), c.msg)
		}
	}
}

func TestProbeConfig_getRunningCtr(t *testing.T) {
	cases := []struct {
		ctr *atomic.Int32
		msg string
	}{
		{
			ctr: new(atomic.Int32),
			msg: "getRunningCtr(ctr) -> new(atomic.Int32)",
		},
		{
			ctr: nil,
			msg: "getRunningCtr(nil) -> new(atomic.Int32)",
		},
	}
	for _, c := range cases {
		cfg := &ProbeConfig{
			RunningCtr: c.ctr,
		}
		if c.ctr != nil {
			assert.Equal(t, c.ctr, cfg.getRunningCtr(), c.msg)
		} else {
			assert.NotNil(t, cfg.getRunningCtr(), c.msg)
		}
	}
}

func TestProbeConfig_getIdleCtr(t *testing.T) {
	cases := []struct {
		ctr *atomic.Int32
		msg string
	}{
		{
			ctr: new(atomic.Int32),
			msg: "getIdleCtr(ctr) -> new(atomic.Int32)",
		},
		{
			ctr: nil,
			msg: "getIdleCtr(nil) -> new(atomic.Int32)",
		},
	}
	for _, c := range cases {
		cfg := &ProbeConfig{
			IdleCtr: c.ctr,
		}
		if c.ctr != nil {
			assert.Equal(t, c.ctr, cfg.getIdleCtr(), c.msg)
		} else {
			assert.NotNil(t, cfg.getIdleCtr(), c.msg)
		}
	}
}

func TestProbeConfig_getWaitGroup(t *testing.T) {
	cases := []struct {
		wg  *sync.WaitGroup
		msg string
	}{
		{
			wg:  new(sync.WaitGroup),
			msg: "getWaitGroup(wg) -> new(sync.WaitGroup)",
		},
		{
			wg:  nil,
			msg: "getWaitGroup(nil) -> new(sync.WaitGroup)",
		},
	}
	for _, c := range cases {
		cfg := &ProbeConfig{
			WaitGroup: c.wg,
		}
		if c.wg != nil {
			assert.Equal(t, c.wg, cfg.getWaitGroup(), c.msg)
		} else {
			assert.NotNil(t, cfg.getWaitGroup(), c.msg)
		}
	}
}
