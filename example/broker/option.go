package broker

import (
	"context"
)

type Options struct {
	// Registry used for clustering TODO
	//Registry registry.Registry

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler
}

type Option func(*Options)

type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type PublishOption func(*PublishOptions)

func PublishWithContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

type SubscribeOptions struct {

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
}

type SubscribeOption func(*SubscribeOptions)
