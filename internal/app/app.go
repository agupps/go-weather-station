package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weather-station/internal/client"
	"weather-station/internal/config"
	"weather-station/internal/metrics"
	"weather-station/internal/models"
	"weather-station/internal/poller"
	"weather-station/internal/weather"

	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type App struct {
	config *config.Config
}

func New() *App {
	c := &config.Config{}
	if err := c.Parse(); err != nil {
		panic(err)
	}
	return &App{config: c}
}

func (a *App) Run() int {

	ctx, cancel := context.WithCancel(context.Background())

	// create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// poller
	p := poller.NewPoller(time.Second, logger)

	// registry + metrics
	registry := prometheus.NewRegistry()
	newMetrics := metrics.NewMetrics(registry)
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// client for open weather api
	apiClient := client.NewOpenWeatherMapClient(a.config.Secret, a.config.WeatherProperties)

	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     a.config.Redis.Addr,     // use default Addr
		Password: a.config.Redis.Password, // no password set
		DB:       a.config.Redis.DB,       // use default DB
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Info("Connected")
			return nil
		},
	})

	redisStore := models.NewRedisStore(redisClient, logger)

	done := a.handleExit(logger)

	go func() {
		defer cancel()
		<-done
		logger.Info("executing cancel()")
	}()

	for _, zipcode := range a.config.Locations {
		w := weather.New(apiClient, zipcode, logger, newMetrics)
		p.Add(w)

		log.Println(zipcode)
		if err := redisStore.Create(ctx, models.NewLocation(zipcode)); err != nil {
			logger.Error("error creating location", "location", zipcode)
			return 1
		}
	}

	// register new "GET /hello" route
	e := echo.New()
	e.GET("/metrics", echo.WrapHandler(promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry})))

	go func() {
		err := e.Start(":8080")
		if err != nil {
			logger.Error("http server hit error", "error", err)
		}
	}()

	p.Start(ctx)
	logger.Info("ending app")

	return 0
}

func (a *App) handleExit(logger *slog.Logger) <-chan struct{} {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)
	done := make(chan struct{})
	go func() {
		<-sig
		logger.Info("Handling exit signal")
		done <- struct{}{}
		close(sig)
		close(done)
	}()
	return done
}
