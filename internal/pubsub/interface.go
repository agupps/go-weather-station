package pubsub

import "context"

type Subscriber[T any] interface {
	Subscribe(context.Context, ...string) chan Event[T]
}

type Publisher[T any] interface {
	Publish(context.Context, string, T) error
}

type PubSub[T any] interface {
	Publisher[T]
	Subscriber[T]
}

type Event[T any] interface {
	Type() string
	Data() T
}
