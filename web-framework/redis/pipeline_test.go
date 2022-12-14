package test

import (
	"context"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.137.100:7000",
		DB:       0,
		PoolSize: 10,
	})
}

func TestPipeline(t *testing.T) {
	ctx := context.Background()
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
