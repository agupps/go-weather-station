package poller

import "time"

type Caller interface {
	Call()
}

type Poller struct {
	pollPeriod time.Duration
	items []Caller
}

func NewPoller(pollPeriod time.Duration) *Poller{
	return &Poller{
		pollPeriod: pollPeriod,
	}
} 

func (p *Poller) Add(item Caller) {
	p.items = append(p.items, item)
}