package gslog

import (
	"log/slog"
	"os"
	"runtime/debug"
	"testing"
)

// 添加功能的属性

func TestGroup(t *testing.T) {
	handler := slog.NewTextHandler(os.Stdout, nil)
	buildInfo, _ := debug.ReadBuildInfo()
	logger := slog.New(handler).With(slog.Group("program", slog.Int("pid", os.Getpid()), slog.String("go_version", buildInfo.GoVersion)))
	slog.SetDefault(logger)
	slog.Info("slog", slog.String("name", "slog"), slog.Int("age", 18))
	// output:
	// time=2024-03-31T17:01:08.875+08:00 level=INFO msg=slog program.pid=32983 program.go_version=go1.22.0 name=slog age=18
}
