package reqpretty

import (
	"context"
	"log/slog"
)

// LogHandler is a slog.Handler implementation that can be used to log to a file.
type LogHandler struct {
	handler slog.Handler
}

// NewHandler creates a new LogHandler.
func NewHandler(logger *Logger) LogHandler {
	if logger == nil {
		logger = Default
	}
	return LogHandler{
		handler: logger,
	}
}

func (l LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return l.handler.Enabled(ctx, level)
}

func (l LogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Customize the logging format here
	return l.handler.Handle(ctx, record)
}

func (l LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LogHandler{
		handler: l.handler.WithAttrs(attrs),
	}
}

func (l LogHandler) WithGroup(name string) slog.Handler {
	return LogHandler{
		handler: l.handler.WithGroup(name),
	}
}
