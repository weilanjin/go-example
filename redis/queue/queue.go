package queue

import (
	"github.com/weilanjin/go-example/redis/initialize"
)

var rdb redis.UniversalClient

func init() {
	rdb = initialize.Redis()
}