package broker

import (
	"fmt"
	"testing"
)

func TestNewRedisPubSub(t *testing.T) {
	mq := NewRedisPubSub(nil)
	var user = struct {
		Name string
	}{
		Name: "test",
	}
	mq.Publish("broker_redis_pubsub_test", &user)
	go mq.Subscribe("broker_redis_pubsub_test", func(e Event) error {
		fmt.Println(e.Message())
		return nil
	})
}
