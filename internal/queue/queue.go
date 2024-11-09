package queue

import "context"

type Queue interface {
	Next(context.Context) (string, error)
	Enqueue(context.Context, string) error
	Dequeue(context.Context) (string, error)
	Delete(context.Context, string) error
}
