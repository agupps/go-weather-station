package queue

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	logger *slog.Logger
}

const (
	queueKey = "queue:locations"
)

func NewRedisQueue(client *redis.Client, logger *slog.Logger) *RedisQueue {
	return &RedisQueue{
		client: client,
		logger: logger,
	}
}

func (r *RedisQueue) Next(ctx context.Context) (string, error) {
	return r.client.LMove(ctx, queueKey, queueKey, "LEFT", "RIGHT").Result()
}

func (r *RedisQueue) Enqueue(ctx context.Context, str string) error {
	if _, err := r.client.RPush(ctx, queueKey, str).Result(); err != nil {
		r.logger.Error("pushing onto queue error", "error", err)
		return err
	}
	return nil
}

func (r *RedisQueue) Dequeue(ctx context.Context) (string, error) {
	str, err := r.client.LPop(ctx, queueKey).Result()
	if err != nil {
		r.logger.Error("popping from queue error", "error", err)
		return str, err
	}
	return str, nil
}

func (r *RedisQueue) Delete(ctx context.Context, str string) error {
	_, err := r.client.LRem(ctx, queueKey, 0, str).Result()
	if err != nil {
		r.logger.Error("error deleting value from queue", "error", err, "value", str)
		return err
	}
	return nil
}
