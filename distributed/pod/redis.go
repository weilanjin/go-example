package pod

import (
	"context"
)

var rdb redis.UniversalClient

func init() {
	// 创建 Redis 客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}