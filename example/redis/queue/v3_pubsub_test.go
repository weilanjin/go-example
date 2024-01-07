package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
	Publisher/Subscriber  发布者订阅者模型
	Redis 消息队列的不足之处，那就是它不支持消息的多播机制。
	消息多播
		允许生产者生产一次消息，中间件负责将消息「复制」到多个消息队列，每个消息队列由相应的消费组进行消费。
	PubSub 会长时间持有连接池的一个网络连接.

	缺点：
		1.Redis 会直接找到相应的消费者传递过去。如果一 个消费者都没有，那么消息直接丢弃。如果开始有三个消费者，
		一个消费者突然挂掉了，生 产者会继续发送消息，另外两个消费者可以持续收到消息。但是挂掉的消费者重新连上
		的时候，这断连期间生产者发送的消息，对于这个消费者来说就是彻底丢失了。
		2.如果 Redis 停机重启，PubSub 的消息是不会持久化的。
*/

func TestPub(t *testing.T) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for i := 1; i < 3600; i++ {
		msg := fmt.Sprintf("this is msg, %d", i)
		if err := rdb.Publish(context.Background(), "queue_v3", msg).Err(); err != nil {
			panic(err)
		}
		fmt.Println("produce:", msg)
		<-ticker.C
	}
}

// 第一个消费者
func TestSub1(t *testing.T) {
	pubSub := rdb.Subscribe(context.Background(), "queue_v3")
	defer pubSub.Close()

	ch := pubSub.Channel()
	for msg := range ch {
		fmt.Println("sub1", msg.String())
	}
}

// 第二个消费者
func TestSub2(t *testing.T) {
	pubSub := rdb.Subscribe(context.Background(), "queue_v3")
	defer pubSub.Close()

	ch := pubSub.Channel()
	for msg := range ch {
		fmt.Println("sub2", msg.String())
	}
}
