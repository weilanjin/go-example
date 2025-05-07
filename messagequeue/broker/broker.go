package broker

type Broker interface {
	Publish(topic string, msg any, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) error
}

type Handler func(Event) error

type Event interface {
	Topic() string
	Message() any
	Ack() error
	Error() error
}
