// Package key 管理

package key

import (
	"context"
	"github.com/redis/go-redis/v9"
	"lovec.wlj/example/redis/initialize"
	"testing"
)

var (
	rdb redis.UniversalClient
	ctx = context.Background()
)

func init() {
	rdb = initialize.Redis()
}

// type
// del
// object
// exists
// expire
// rename 如果 newkey 已经存在那么它的值也将被覆盖 -- renamenx 更安全 存在就重命名不成功。
// randomkey 随机返回一个 key
func TestKey(t *testing.T) {
	rdb.Set(ctx, "golang", "1.22", 0)

	// rename 如果 go 已经存在那么它的值也将被复制
	//
	// renamenx 存在则返回0
	// 重命名期间会删除 oldkey， 如果旧的 key 大会导致 redis 阻塞
	rdb.RenameNX(ctx, "golang", "go")
}
