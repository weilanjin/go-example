package debug

import (
	internalTrace "lovec.wlj/example/debug/internal/trace"
	"lovec.wlj/example/debug/trace"
)

func Trace(traceID, span, tip string, opts ...trace.Option) {
	internalTrace.Debug(traceID, span, tip, func(o *internalTrace.TraceLog) {
		for _, opt := range opts {
			opt(o)
		}
	})
}
