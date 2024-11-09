package poller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	client *redis.Client
	uid    uuid.UUID
	logger *slog.Logger
}

func NewRedisLock(client *redis.Client, logger *slog.Logger) *RedisLock {
	return &RedisLock{
		client: client,
		uid:    uuid.New(),
		logger: logger,
	}
}

const (
	lockKey     = "poller:lock"
	lockTimeout = 5 * time.Second
)

func (l *RedisLock) Lock(ctx context.Context) bool {
	str, err := l.client.Get(ctx, lockKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			set, err := l.client.SetNX(ctx, lockKey, l.uid.String(), lockTimeout).Result()
			if err != nil {
				return false
			}
			if set {
				l.logger.Info("lock obtained", "uid", l.uid.String())
				return true
			}
		} else if err != nil {
			l.logger.Error("redis client hit err", "error", err)
			return false
		}
	}
	l.logger.Info("lock is held", "current lock id", str, "current thread id", l.uid.String())
	if str == l.uid.String() {
		if _, err := l.client.Expire(ctx, lockKey, lockTimeout).Result(); err != nil {
			l.logger.Error("error updating lock expiration", "error", err)
			return false
		}
		l.logger.Info("lock retained")
		return true
	}
	l.logger.Info("lock is held by different thread")
	return false
}
