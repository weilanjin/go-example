package main

import (
	"log/slog"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func TestSlogZap(t *testing.T) {
	zapL, _ := zap.NewDevelopment()
	defer zapL.Sync()

	slog.SetDefault(slog.New(zapslog.NewHandler(zapL.Core())))

	slog.Info("sample log message", slog.String("field1", "value1"), slog.Int("field2", 33))
}
