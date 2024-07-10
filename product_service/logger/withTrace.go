package logger

import (
	"context"
	"io"
	"log/slog"
)

type WithTraceID struct {
	h slog.Handler
}

func NewTraceHandler(w io.Writer, opts *slog.HandlerOptions) *WithTraceID {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	return &WithTraceID{
		h: slog.NewJSONHandler(w, opts),
	}
}

func (t *WithTraceID) Enabled(context.Context, slog.Level) bool {
	return true
}

func (t *WithTraceID) Handle(ctx context.Context, r slog.Record) error {

	return t.h.Handle(ctx, r)
}

func (t *WithTraceID) WithAttrs(attrs []slog.Attr) slog.Handler {
	return t.h.WithAttrs(attrs)
}

func (t *WithTraceID) WithGroup(name string) slog.Handler {
	return t.h.WithGroup(name)
}
