package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"runtime/debug"

	"github.com/redis/go-redis/v9"
)

type redisPubSub struct {
	rdb  redis.UniversalClient
	opts *Options
}

func NewRedisPubSub(rdb redis.UniversalClient, opts ...Option) Broker {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	return &redisPubSub{
		rdb:  rdb,
		opts: &options,
	}
}

func (b *redisPubSub) Publish(topic string, msg any, opts ...PublishOption) error {
	var options PublishOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.Context == nil {
		options.Context = context.Background()
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return b.rdb.Publish(options.Context, topic, bytes).Err()
}

func (b *redisPubSub) Subscribe(topic string, h Handler, opts ...SubscribeOption) error {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			slog.Error("Subscribe panic:%v", slog.Any("err", err))
		}
	}()

	var options SubscribeOptions
	for _, o := range opts {
		o(&options)
	}
	subscribe := b.rdb.Subscribe(options.Context, topic)
	for {
		select {
		case msg := <-subscribe.Channel():
			event := &redisPubSubEvent{
				message: msg.Payload,
				topic:   msg.Channel,
			}
			event.err = h(event)
			if event.err != nil && b.opts.ErrorHandler != nil {
				_ = b.opts.ErrorHandler(event)
			}
		case <-options.Context.Done():
			return nil
		}
	}
}

type redisPubSubEvent struct {
	err     error
	message any
	opts    *Options
	topic   string
}

func (r *redisPubSubEvent) Topic() string {
	return r.topic
}

func (r *redisPubSubEvent) Message() any {
	return r.message
}

func (r *redisPubSubEvent) Ack() error {
	return nil
}

func (r *redisPubSubEvent) Error() error {
	return r.err
}
