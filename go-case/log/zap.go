package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zapLogger *zap.Logger
}

func NewZap(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config := zap.Config{
		Level: zap.NewAtomicLevelAt(zapLevel),
	}
	l, err := config.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(err)
	}
	return &zapLogger{
		zapLogger: l.Named(opts.Name),
	}
}

func Init(opts *Options) {
	mu.RLock()
	defer mu.RUnlock()
	defaultLogger = NewZap(opts)
}

func (z *zapLogger) Debug(msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Debugf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Debugw(msg string, keysAndValues ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Debugc(ctx context.Context, msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Debugcf(ctx context.Context, format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Debugcw(ctx context.Context, msg string, keysAndValues ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Eorror(msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Eorrorf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Eorrorw(msg string, keysAndValues ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Eorrorc(ctx context.Context, msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Eorrorcf(ctx context.Context, format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Eorrorcw(ctx context.Context, msg string, keysAndValues ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Info(msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Infof(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Infow(msg string, keysAndValues ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Infoc(ctx context.Context, msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Infocf(ctx context.Context, format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Infocw(ctx context.Context, msg string, keysAndValues ...any) {
	//TODO implement me
	panic("implement me")
}