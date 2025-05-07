package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
	Stream
		Redis v5.0
	1.强大的支持多播的可持久化的消息队列.
	2.每个消息都有一个唯一的ID和对应的内容.
	3.消息是持久化的，Redis 重启后，内容还在.
*/

func TestStreamProduce(t *testing.T) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	ctx := context.Background()

	for i := 1; i < 3600; i++ {
		msg := fmt.Sprintf("this is msg, %d", i)
		if err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream:     "queue_v4",                                     // stream key
			NoMkStream: false,                                          // * 默认false,当为false时,key不存在，会新建
			MaxLen:     100000,                                         // * 指定stream的最大长度,当队列长度超过上限后，旧消息会被删除，只保留固定长度的新消息
			Approx:     false,                                          // * 默认false,当为true时,模糊指定stream的长度
			ID:         "*",                                            // 消息 id，我们使用 * 表示由 redis 生成
			Values:     []interface{}{"content", msg, "name", "lance"}, // 表示消息内容键值对
			// MinID: "id",            // * 超过阈值，丢弃设置的小于MinID消息id【基本不用】
			// Limit: 1000,            // * 限制长度【基本不用】
		}).Err(); err != nil {
			fmt.Println(err)
		}

		fmt.Println("produce:", msg)
		<-ticker.C
	}
}

// 第一个消费者
func TestStreamConsume1(t *testing.T) {
	ctx := context.Background()
	if err := rdb.XGroupCreate(ctx, "queue_v4", "group1", "$").Err(); err != nil { // 0:从头获取   $：从最新获取
		if !redis.HasErrorPrefix(err, "BUSYGROUP Consumer Group name already exists") {
			panic(err)
		}
	}

	for {
		r, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "group1",
			Consumer: "consumer1",
			Streams:  []string{"queue_v4", ">"},
			Count:    1,
			Block:    0,
			NoAck:    true, // 为true，表示读取消息时确认消息
		}).Result()
		// id为0或其他，表示可获取已读但未确认的消息，请注意，在这种情况下，BLOCK和NoACK都将被忽略。
		// id为>，表示着消费者只希望接收从未传递给任何其他消费者的消息
		if err != nil {
			panic(err)
		}
		for _, v := range r {
			fmt.Println(v.Messages)
		}
	}
}

// 第二个消费者
func TestStreamConsume2(t *testing.T) {
	ctx := context.Background()
	if err := rdb.XGroupCreate(ctx, "queue_v4", "group2", "0").Err(); err != nil { // 0:从头获取   $：从最新获取
		if !redis.HasErrorPrefix(err, "BUSYGROUP Consumer Group name already exists") {
			panic(err)
		}
	}

	for {
		r, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "group1",
			Consumer: "consumer1",
			Streams:  []string{"queue_v4", ">"},
			Count:    1,
			Block:    0,
			NoAck:    true, // 为true，表示读取消息时确认消息
		}).Result()
		// id为0或其他，表示可获取已读但未确认的消息，请注意，在这种情况下，BLOCK和NoACK都将被忽略。
		// id为>，表示着消费者只希望接收从未传递给任何其他消费者的消息
		if err != nil {
			panic(err)
		}
		for _, v := range r {
			fmt.Println(v.Messages)
		}
	}
}