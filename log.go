package probe

import (
	"context"
	"log/slog"
)

type (
	// NoopLogHandler is a no-op log handler.
	NoopLogHandler struct{}
)

// Enabled implementation of slog.Handler for probe.NoopLogHandler.
// Returns false for all log levels.
func (h *NoopLogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle implementation of slog.Handler for NoopLogHandler.
func (h *NoopLogHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs implementation of slog.Handler for NoopLogHandler.
func (h *NoopLogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup implementation of slog.Handler for NoopLogHandler.
func (h *NoopLogHandler) WithGroup(_ string) slog.Handler {
	return h
}
