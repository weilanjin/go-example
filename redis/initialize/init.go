package initialize

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Redis() redis.UniversalClient {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		ClientName: "example",
	})
	ctx := context.Background()
	ping := rdb.Ping(ctx)
	if ping.Err() != nil {
		panic(ping.Err())
	}
	log.Println(ping)
	log.Println(rdb.DBSize(ctx))
	return rdb
}
