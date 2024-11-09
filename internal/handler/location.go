package handler

import (
	"context"
	"log/slog"
	"weather-station/internal/models"
	"weather-station/internal/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

type Location struct {
	model  models.Model
	logger *slog.Logger
}

func NewLocation(model models.Model, logger *slog.Logger) *Location {
	return &Location{
		model:  model,
		logger: logger,
	}
}

func (l *Location) List(c echo.Context) error {
	ctx := context.Background()

	locations, err := l.model.List(ctx)
	if err != nil {
		l.logger.Error("error getting all locations", "error", err)
		return err
	}

	component := templates.Location(locations)

	h := templ.Handler(component)

	return h.Component.Render(ctx, c.Response().Writer)
}
