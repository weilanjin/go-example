package queue

import (
	"github.com/redis/go-redis/v9"
	"github.com/weilanjin/go-example/microservice/redis/initialize"
)

var rdb redis.UniversalClient

func init() {
	rdb = initialize.Redis()
}