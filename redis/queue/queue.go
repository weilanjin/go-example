package queue

import (
	"github.com/redis/go-redis/v9"
	"github.com/weilanjin/go-example/redis/initialize"
)

var rdb redis.UniversalClient

func init() {
	rdb = initialize.Redis()
}
