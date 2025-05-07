package cmdscan

import (
	"context"
	"log"
	"log/slog"
	"testing"

	"github.com/weilanjin/go-example/redis/initialize"
)

var (
	rdb redis.UniversalClient
	ctx = context.Background()
)

func init() {
	rdb = initialize.Redis()
}

func TestScan(t *testing.T) {
	var total int
	var (
		keys []string
		cur  uint64
		err  error
	)
	for {
		keys, cur, err = rdb.Scan(ctx, cur, "*", 2).Result()
		if err != nil {
			slog.Error(err.Error())
			return
		}
		total += len(keys)

		slog.Info("result ", "keys", keys, "cur", cur)
		if cur == 0 {
			break
		}
	}

	log.Println("count", total)
}