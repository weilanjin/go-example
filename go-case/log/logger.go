package logger

import (
	"context"
	"sync"
)

type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...any)
	Debugw(msg string, keysAndValues ...any)
	Debugc(ctx context.Context, msg string)
	Debugcf(ctx context.Context, format string, args ...any)
	Debugcw(ctx context.Context, msg string, keysAndValues ...any)

	Eorror(msg string)
	Eorrorf(format string, args ...any)
	Eorrorw(msg string, keysAndValues ...any)
	Eorrorc(ctx context.Context, msg string)
	Eorrorcf(ctx context.Context, format string, args ...any)
	Eorrorcw(ctx context.Context, msg string, keysAndValues ...any)

	Info(msg string)
	Infof(format string, args ...any)
	Infow(msg string, keysAndValues ...any)
	Infoc(ctx context.Context, msg string)
	Infocf(ctx context.Context, format string, args ...any)
	Infocw(ctx context.Context, msg string, keysAndValues ...any)
}

var (
	defaultLogger = NewZap(NewOptions())
	mu            sync.RWMutex
)