package gslog

import (
	"context"
	slogcontext "github.com/PumpkinSeed/slog-context"
	"github.com/phsym/console-slog"
	slogmulti "github.com/samber/slog-multi"
	slogsampling "github.com/samber/slog-sampling"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestMultiHandler(t *testing.T) {
	// Will print 33% of entries.
	sampling := slogsampling.UniformSamplingOption{
		// The sample rate for sampling traces in the range [0.0, 1.0].
		Rate: 0.33,
	}

	out := console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:      slog.LevelDebug,
		TimeFormat: time.DateTime + ".000",
		AddSource:  true,
	})
	slogHandler := slogmulti.
		Pipe(slogcontext.NewHandler).
		Pipe(sampling.NewMiddleware()).
		Handler(out)

	slog.SetDefault(slog.New(slogHandler))
	//slog.New(slogmulti.Fanout())
	ctx := slogcontext.WithValue(context.Background(), "slogcontext", "slogcontext value")
	for i := 0; i < 15; i++ {
		slog.InfoContext(ctx, "hello", slog.Int("idx", i), slog.String("name", "samber"))
	}
}
