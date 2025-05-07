package debug

import (
	internalTrace "github.com/weilanjin/go-example/debug/internal/trace"
	"github.com/weilanjin/go-example/debug/trace"
)

func Trace(traceID, span, tip string, opts ...trace.Option) {
	internalTrace.Debug(traceID, span, tip, func(o *internalTrace.TraceLog) {
		for _, opt := range opts {
			opt(o)
		}
	})
}