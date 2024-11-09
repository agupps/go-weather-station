package pubsub

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type RedisEvent struct {
	eventType, eventData string
}

func NewRedisEvent(eventType, eventData string) RedisEvent {
	return RedisEvent{
		eventType: eventType,
		eventData: eventData,
	}
}

func (re RedisEvent) Type() string { return re.eventType }
func (re RedisEvent) Data() string { return re.eventData }

type RedisPubSub struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisPubSub(client *redis.Client, logger *slog.Logger) *RedisPubSub {
	return &RedisPubSub{
		client: client,
		logger: logger,
	}
}

func (ps *RedisPubSub) Publish(ctx context.Context, event string, data string) error {
	if _, err := ps.client.Publish(ctx, event, data).Result(); err != nil {
		return err
	}
	return nil
}

func (ps *RedisPubSub) Subscribe(ctx context.Context, events ...string) chan Event[string] {
	out := make(chan Event[string])
	sub := ps.client.Subscribe(ctx, events...)
	dataChannel := sub.Channel()
	go func() {
		ps.logger.Info("receiving from subscription")
		for {
			select {
			case data := <-dataChannel:
				ps.logger.Info("receiving data", "event", data.Channel, "data", data.Payload)
				out <- NewRedisEvent(data.Channel, data.Payload)
			case <-ctx.Done():
				ps.logger.Info("receiving done signal from subscriber")
				close(out)
				return
			}
		}
	}()
	return out
}
