package models

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

const locationsSet = "Locations"

type RedisStore struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisStore(redisClient *redis.Client, logger *slog.Logger) *RedisStore {
	return &RedisStore{
		client: redisClient,
		logger: logger.With("context", "redis_store"),
	}
}

func (r *RedisStore) Create(ctx context.Context, location *Location) error {
	_, err := r.client.SAdd(ctx, locationsSet, location.Zipcode).Result()
	if err != nil {
		return err
	}
	jsonrep, _ := json.Marshal(location)
	if _, err := r.client.Set(ctx, location.Zipcode, jsonrep, 0).Result(); err != nil {
		return err
	}

	return nil // no error
}

func (r *RedisStore) Update(ctx context.Context, location *Location) error {
	return r.Create(ctx, location)
}

func (r *RedisStore) Delete(ctx context.Context, location *Location) error {
	count, err := r.client.SRem(ctx, locationsSet, location.Zipcode).Result()
	if err != nil {
		return nil
	}
	if count > 0 {
		if _, err := r.client.Del(ctx, location.Zipcode).Result(); err != nil {
			return err
		}
	}
	return nil // no error
}

func (r *RedisStore) List(ctx context.Context) ([]Location, error) {
	zipCodes, err := r.client.SMembers(ctx, locationsSet).Result()
	if err != nil {
		return nil, err
	}
	var locations []Location
	for _, zip := range zipCodes {
		location, err := r.Get(ctx, zip)
		if err != nil {
			r.logger.Error("Error getting location info for zipcode", "zipcode", zip, "error", err)
			continue
		}
		locations = append(locations, *location)
	}
	return locations, nil
}

func (r *RedisStore) Get(ctx context.Context, zipCode string) (*Location, error) {
	str, err := r.client.Get(ctx, zipCode).Result()
	if err != nil {
		return nil, err
	}
	location := &Location{}
	if err := json.Unmarshal([]byte(str), location); err != nil {
		return nil, err
	}
	return location, nil
}

func (r *RedisStore) Notify(l *Location) {
	r.logger.Info("redis store calling update", "location", l)

	if err := r.Update(context.TODO(), l); err != nil {
		r.logger.Error("error updating db", "error", err)
	}
	return
}
