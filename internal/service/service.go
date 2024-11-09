package service

import "context"

type Service interface {
	Start(ctx context.Context)
}
