package gslog

import (
	"log"
	"log/slog"
	"os"
	"testing"
)

func TestSetDefault(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       nil,
		ReplaceAttr: nil,
	}))
	// SetDefault 也会影响 log.Println
	slog.SetDefault(logger)
	slog.Info("test slog info")
	log.Println("test log info")
	// Output:
	// time=2024-03-31T16:44:01.930+08:00 level=INFO source=/Users/lanjin/Documents/work/code/go-example/go-sourcecode/slog/slog_set_default_test.go:17 msg="test slog info"
	// time=2024-03-31T16:44:01.930+08:00 level=INFO source=:0 msg="test log info"
}
