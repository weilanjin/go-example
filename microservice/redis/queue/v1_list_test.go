package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
	queue v1
		rpush + lpop / lpush + rpop
		rpush: 在列表中添加一个或多个值到列表尾部
		lpop:  移出并获取列表的第一个元素

	1. 如果队列空了，客户端就会陷入 pop 的死循环,空轮询不但拉高了客户端的 CPU，redis 的 QPS 也 会被拉高.
	2. sleep 来解决这个问题, CPU 能降下来,QPS 也降下来了.
	3. 睡眠会导致消息的延迟增大. 阻塞读解决 blpop/brpop -- queue v2
*/

func TestListQueueProduce(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for i := 1; i < 3600; i++ {
		msg := fmt.Sprintf("this is msg, %d", i)
		if err := rdb.RPush(context.Background(), "queue_v1", msg).Err(); err != nil {
			panic(err)
		}
		fmt.Println("produce:", msg)
		<-ticker.C
	}
}

func TestListQueueConsume(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		msg, err := rdb.LPop(context.Background(), "queue_v1").Result()
		if err != nil {
			fmt.Println("err:", err)
		} else {
			fmt.Println("consume:", msg)
		}
		<-ticker.C
	}
}

func TestMget(t *testing.T) {
	result, err := rdb.MGet(context.Background(), "go2", "go1").Result()
	fmt.Println("result:", result)
	fmt.Println("err:", err)
}
