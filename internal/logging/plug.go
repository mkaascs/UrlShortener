package logging

import (
	"context"
	"log/slog"
)

func NewPlugLogger() *slog.Logger {
	return slog.New(NewPlugHandler())
}

type PlugHandler struct{}

func NewPlugHandler() *PlugHandler {
	return &PlugHandler{}
}

func (h *PlugHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (h *PlugHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (h *PlugHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PlugHandler) WithGroup(name string) slog.Handler {
	return h
}
