package gslog

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

// 内置日志级别
// DEBUG -4
// INFO -0
// WARN 4
// ERROR -8

// 自定义日志级别

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Level]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func SetLevelName(a slog.Attr) slog.Attr {
	level := a.Value.Any().(slog.Level)
	levelLabel, exists := LevelNames[level]
	if !exists {
		levelLabel = level.String()
	}
	a.Value = slog.StringValue(levelLabel)
	return a
}

func TestCustomLevel(t *testing.T) {
	logLevel := &slog.LevelVar{}
	opts := slog.HandlerOptions{
		Level: logLevel, // 默认INFO
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				a = SetLevelName(a)
			}
			return a
		},
	}
	logLevel.Set(LevelTrace) // 动态过滤级别

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &opts)))
	ctx := context.Background()
	slog.Log(ctx, LevelTrace, "trace")
	slog.Log(ctx, LevelFatal, "fatal")
}
