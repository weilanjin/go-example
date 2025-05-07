package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"log/slog"
	"testing"
)

func TestSlogZap(t *testing.T) {
	zapL, _ := zap.NewDevelopment()
	defer zapL.Sync()

	slog.SetDefault(slog.New(zapslog.NewHandler(zapL.Core(), &zapslog.HandlerOptions{
		AddSource: true,
	})))

	slog.Info("sample log message", slog.String("field1", "value1"), slog.Int("field2", 33))
}
