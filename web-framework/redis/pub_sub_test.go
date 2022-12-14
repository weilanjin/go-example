package test

import (
	"context"
	"log"
	"testing"
)

// Redis 发布的消息不会被持久化，这就会导致新订阅的客户端将不会收到历史消息

func TestPub(t *testing.T) {
	ctx := context.Background()
	if err := rdb.Publish(ctx, "mychannel", "payload").Err(); err != nil {
		t.Fatal(err)
	}
}

// debug pattern
func TestSub(t *testing.T) {
	ctx := context.Background()
	ps := rdb.Subscribe(ctx, "mychannel")
	defer ps.Close()
	for {
		m, err := ps.ReceiveMessage(ctx)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(m.Channel, m.Payload)
	}
}
