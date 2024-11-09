package poller

import "context"

type Locker interface {
	Lock(context.Context) bool
}
