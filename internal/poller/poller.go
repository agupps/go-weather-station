package poller

import (
	"log/slog"
	"time"
)

type Caller interface {
	Call()
}

type Poller struct {
	pollPeriod time.Duration
	items []Caller
	logger *slog.Logger
}

func NewPoller(pollPeriod time.Duration, logger *slog.Logger) *Poller{
	return &Poller{
		pollPeriod: pollPeriod,
		logger: logger,
	}
} 

func (p *Poller) Add(item Caller) {
	p.items = append(p.items, item)
}

func (p *Poller) Start() {
	ticker := time.NewTicker(p.pollPeriod)
	counter := 0
	for {
		<- ticker.C
		index := counter % len(p.items)

		go p.items[index].Call()
		counter++
	} 
}