package gslog

import (
	"github.com/weilanjin/go-example/log/slog/deep"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// 写入日志文件并输出到控制台
func TestDefault(t *testing.T) {
	lvl := &slog.LevelVar{}
	lvl.Set(slog.LevelError)

	writerFile := readFile("./foo.log")
	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stderr, writerFile), &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// set time format.
			if a.Key == slog.TimeKey && len(groups) == 0 {
				format := a.Value.Time().Format(time.DateTime + ".000")
				return slog.String(slog.TimeKey, format)
			}
			// Remove the directory from the source's filename.
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}))
	slog.SetDefault(logger)
	slog.Debug("Debug level logging", "level", "debug")
	slog.Info("info level logging", slog.String("level", "info"))
	slog.Warn("Warn level logging", "level", "warn")
	slog.Error("Error level logging", "level", "error")
	lvl.Set(slog.LevelDebug)
	slog.Debug("Debug level logging")
	deep.Add(12, 14)

	// for i := 0; i < 100000; i++ {
	//	slog.Info("greeting", "say", "hello")
	// }
	// output:
	// {"time":"2023-09-11 09:40:10.806","level":"INFO","source":{"function":"lovec.wlj/go-sourcecode/slog.TestDefault","file":"default_test.go","line":38},"msg":"greeting","say":"hello"}
}

func readFile(filename string) io.Writer {
	r := &lumberjack.Logger{
		Filename:   filename,
		LocalTime:  true,
		MaxSize:    1,
		MaxAge:     3,
		MaxBackups: 5,
		Compress:   true,
	}
	return r
}