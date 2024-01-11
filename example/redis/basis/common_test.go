package basis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"lovec.wlj/example/redis/initialize"
	"testing"
)

var rdb redis.UniversalClient

func init() {
	rdb = initialize.Redis()
}

func TestCommon(t *testing.T) {
	ctx := context.Background()

	// dbsize 命令在计算键总数时不会遍历所有键，而是直接获取Redis内置的 键总数变量，所以dbsize命令的时间复杂度是O（1）
	rdb.DBSize(ctx)
}
