package logging

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
This file tests The NoopLogHandler implementation of the slog.Handler interface.
The implementation is purpose built to do nothing, so these tests may be somewhat silly.
However, they are included for completeness and test coverage.
*/

func TestNoopLogHandler_Enabled(t *testing.T) {
	cases := []struct {
		l   slog.Level
		msg string
	}{
		{
			l:   slog.LevelDebug,
			msg: "Enabled(debug) -> false",
		},
		{
			l:   slog.LevelInfo,
			msg: "Enabled(info) -> false",
		},
		{
			l:   slog.LevelWarn,
			msg: "Enabled(warn) -> false",
		},
		{
			l:   slog.LevelError,
			msg: "Enabled(error) -> false",
		},
	}
	h := &NoopLogHandler{}
	ctx := context.Background()
	for _, c := range cases {
		assert.False(t, h.Enabled(ctx, c.l), c.msg)
	}
}

func TestNoopLogHandler_Handle(t *testing.T) {
	h := &NoopLogHandler{}
	ctx := context.Background()
	r := slog.Record{}
	assert.NoError(t, h.Handle(ctx, r))
}

func TestNoopLogHandler_WithAttrs(t *testing.T) {
	h := &NoopLogHandler{}
	attrs := []slog.Attr{}
	assert.Equal(t, h, h.WithAttrs(attrs))
}

func TestNoopLogHandler_WithGroup(t *testing.T) {
	h := &NoopLogHandler{}
	group := "test"
	assert.Equal(t, h, h.WithGroup(group))
}
