package trace

import (
	"context"
	"lovec.wlj/example/debug/internal/trace"
)

type TraceLog = trace.TraceLog

func Use(ctx context.Context, c trace.Connector) {
	trace.Use(ctx, c)
}

type Option func(*trace.TraceLog)

type KV map[string]any

func WithData(data any) Option {
	return func(o *trace.TraceLog) {
		o.LogData = data
	}
}

func WithErr(err error) Option {
	return func(o *trace.TraceLog) {
		if err == nil {
			return
		}
		o.Err = err.Error()
	}
}
