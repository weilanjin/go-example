package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
	queue v2
		rpush + blpop / lpush + brpop
		1. blocking 阻塞读在队列没有数据的时候，会立即进入休眠状态，一旦数据到来，则立刻醒过来。消 息的延迟几乎为零。
		2. 需要设置合理的超时时间。
	Redis 的客户端连接就成了闲置连接，闲置过久，服务器一般 会主动断开连接，减少闲置资源占用。这个时候 blpop/brpop 会抛出异常来
*/

// first go test TestListQueueProduce
func TestBlockListQueueConsume(t *testing.T) {
	for {
		msgArr, err := rdb.BLPop(context.Background(), 10*time.Second, "queue_v1").Result()
		if err != nil {
			fmt.Println("err:", err)
		} else {
			fmt.Println("consume:", msgArr)
		}
	}
}
