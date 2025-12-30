package manager

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/weilanjin/go-example/microservice/redis/initialize"
)

var (
	rdb redis.UniversalClient
	ctx = context.Background()
)

func init() {
	rdb = initialize.Redis()
}

func TestSlowLog(t *testing.T) {

	// slowlog-log-slower-than = 10000 us (微妙) = 10ms【 =0 记录所有的命令，<0 不会记录 】
	// slowlog-max-len 慢查询日志最多存储多少条

	// config set slowlog-log-slower-than 20000
	// config set slowlog-max-len 1000
	// config rewrite 将配置持久化到本地配置文件
	rdb.ConfigSet(ctx, "slowlog-log-slower-than", "0")
	slowLogs, _ := rdb.SlowLogGet(ctx, 10).Result()
	for _, sl := range slowLogs {
		// 慢查询只记录命令执行时间
		// 慢查询日志是一个先进先出的队列
		// 可能会丢失部分慢查询命令
		log.Println("慢查询日志标识ID:", sl.ID, " 发生时间戳:", sl.Time, " 命令耗时:", sl.Duration, " 执行命令和参数:", sl.Args)
	}

	// slowlog len 获取慢查询日志列表当前的长度
	// slowlog reset 慢查询日志重置 (对列表做清理操作)
}
