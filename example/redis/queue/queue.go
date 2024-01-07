package queue

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var rdb redis.UniversalClient

func init() {
	rdb = redis.NewUniversalClient(&redis.UniversalOptions{})
	ping, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(ping)
}

type queue interface {
	Produce(ctx context.Context, key string, msg any) error
	Consume(ctx context.Context, key string) error
}
