package gslog_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"
	"time"

	gslog "github.com/weilanjin/go-example/log/slog"
)

func TestCustomHandler(t *testing.T) {
	ch := make(chan []byte, 100)
	attrs := []slog.Attr{
		{Key: "foo1", Value: slog.StringValue("bar1")},
		{Key: "foo2", Value: slog.StringValue("bar2")},
	}
	slog.SetDefault(slog.New(gslog.NewChannelHandler(ch, &slog.HandlerOptions{}).WithAttrs(attrs)))
	go func() { // 模拟channel的消费者，用来消费日志
		for {
			b := <-ch
			fmt.Println(string(b))
		}
	}()
	slog.Info("hello", "name", "Al")
	slog.Error("oops", "err", net.ErrClosed, "status", 500)
	slog.LogAttrs(context.Background(), slog.LevelError, "oops",
		slog.Int("status", 500), slog.Any("err", net.ErrClosed),
	)

	time.Sleep(3 * time.Second)
}
