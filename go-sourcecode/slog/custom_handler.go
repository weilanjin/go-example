package gslog

import (
	"bytes"
	"context"
	"log/slog"
)

type ChannelHandler struct {
	slog.Handler
	ch  chan []byte
	buf *bytes.Buffer
}

func NewChannelHandler(ch chan []byte, opts *slog.HandlerOptions) *ChannelHandler {
	var b = make([]byte, 256)
	h := &ChannelHandler{
		buf: bytes.NewBuffer(b),
		ch:  ch,
	}
	h.Handler = slog.NewTextHandler(h.buf, opts)
	return h
}

func (h *ChannelHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *ChannelHandler) Handle(ctx context.Context, r slog.Record) error {
	if err := h.Handler.Handle(ctx, r); err != nil {
		return err
	}
	buf := make([]byte, h.buf.Len())
	copy(buf, h.buf.Bytes())
	h.ch <- buf
	h.buf.Reset()
	return nil
}

func (h *ChannelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ChannelHandler{
		ch:      h.ch,
		buf:     h.buf,
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *ChannelHandler) WithGroup(name string) slog.Handler {
	return &ChannelHandler{
		ch:      h.ch,
		buf:     h.buf,
		Handler: h.Handler.WithGroup(name),
	}
}