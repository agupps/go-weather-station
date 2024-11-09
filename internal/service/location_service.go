package service

import (
	"context"
	"log/slog"
	"weather-station/internal/client"
	"weather-station/internal/metrics"
	"weather-station/internal/models"
)

type LocationService struct {
	client  client.HTTPGetter
	model   models.Model
	logger  *slog.Logger
	metrics metrics.Observable
}

func NewLocationService(client client.HTTPGetter, model models.Model,
	logger *slog.Logger, metrics metrics.Observable) *LocationService {
	return &LocationService{
		client:  client,
		model:   model,
		logger:  logger,
		metrics: metrics,
	}
}

func (s *LocationService) Start(ctx context.Context) {
	return
}
