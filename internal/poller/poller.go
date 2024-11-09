package poller

import (
	"context"
	"log/slog"
	"time"
	"weather-station/internal/models"
	"weather-station/internal/queue"
)

type Caller interface {
	Call()
}

type Poller struct {
	pollPeriod time.Duration
	logger     *slog.Logger
	queue      queue.Queue
	locker     Locker
}

func NewPoller(pollPeriod time.Duration, logger *slog.Logger, queue queue.Queue, locker Locker) *Poller {
	return &Poller{
		pollPeriod: pollPeriod,
		logger:     logger,
		queue:      queue,
		locker:     locker,
	}
}

func (p *Poller) Add(ctx context.Context, item string) {
	err := p.queue.Enqueue(ctx, item)
	if err != nil {
		p.logger.Error("received error enqueuing", "item", item, "error", err)
		return
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.pollPeriod)
	for {
		select {
		case <-ticker.C:
			if p.locker.Lock(ctx) {
				zip, err := p.queue.Next(ctx)
				if err != nil {
					p.logger.Error("error getting next zip from queue", "error", err)
					continue
				}
				loc := models.NewLocation(zip)

			}

			// ToDo: Actually make the API call and then update datastore
		case <-ctx.Done():
			ticker.Stop()
			p.logger.Info("poller stopping")
			return
		}
	}
}
