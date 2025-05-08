package pipe

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/weilanjin/go-example/redis/initialize"
)

var rdb redis.UniversalClient

func init() {
	rdb = initialize.Redis()
}

func TestPipeline(t *testing.T) {
	ctx := context.Background()
	// rdb.MGet(ctx, "username", "age", "address", "python", "xxx")
	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Get(ctx, "username")
		pipe.Get(ctx, "age")
		pipe.Get(ctx, "address")
		pipe.Get(ctx, "python")
		pipe.Get(ctx, "xxx")
		return nil
	})

	// 忽略没有找到的错误
	if err != nil && err != redis.Nil {
		t.Fatal(err)
	}
	for _, cmd := range cmds {
		log.Println(cmd.(*redis.StringCmd).Val())
	}
}
