package basis

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/weilanjin/go-example/redis/initialize"
)

var (
	rdb redis.UniversalClient
	ctx = context.Background()
)

func init() {
	rdb = initialize.Redis()
}

func TestCommon(t *testing.T) {
	// dbsize 命令在计算键总数时不会遍历所有键，而是直接获取Redis内置的 键总数变量，所以dbsize命令的时间复杂度是O（1）
	rdb.DBSize(ctx)

	// =================================================================
	// 发送命令、执行命令、返回结果
	// 所有命令都会进入一个队列中，然后逐个执行
	// ================================================================

	// 很多存储系统和编程语言内部使用CAS机制实现计数功能，会有一定的 CPU开销，但在Redis中完全不存在这个问题，
	// 因为Redis是单线程架构，任 何命令到了Redis服务端都要顺序执行。
	//
	// 不是 int 会返回错误结果 ERR value is not an integer or out of range
	// incr incrby incrbyfloat
	// decr decrby
	log.Println(rdb.Incr(ctx, "token:limit"))
	log.Println(rdb.IncrBy(ctx, "token:limit", 3))
}
