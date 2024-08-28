package models

import "context"

type Model interface {
	Create(context.Context, *Location) error
	Update(context.Context, *Location) error
	Delete(context.Context, *Location) error
	Get(context.Context, string) (*Location, error)
	List(context.Context) ([]Location, error)
}
