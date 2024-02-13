package main

import (
	"errors"
	"log"
	"log/slog"
	"weather-station/internal/models"

	"net/http"
	"os"
	"time"
	"weather-station/internal/client"
	"weather-station/internal/config"
	"weather-station/internal/metrics"
	"weather-station/internal/poller"
	"weather-station/internal/weather"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	app := pocketbase.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	p := poller.NewPoller(time.Second, logger)
	c := &config.Config{}
	if err := c.Parse(); err != nil {
		panic(err)
	}

	registry := prometheus.NewRegistry()

	newMetrics := metrics.NewMetrics(registry)

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	client := client.NewOpenWeatherMapClient(c.Secret, c.WeatherProperties)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		locations := models.NewLocations(app.Dao())
		for _, zipcode := range c.Locations {
			w := weather.New(client, zipcode, logger, newMetrics)
			p.Add(w)
			w.AddSubscriber(locations)

			log.Println(zipcode)
			if err := locations.Create(models.NewLocation(zipcode)); err != nil {
				return err
			}
		}
		return nil
	})

	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		locations := models.NewLocations(app.Dao())
		if err := locations.Delete(); err != nil {
			return err
		}

		log.Println("terminating...")
		return nil
	})

	go p.Start()

	// Expose /newMetrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	go func() {
		if err := http.ListenAndServe(":8081", nil); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error starting or shutting down HTTP server", "err", err)
			os.Exit(1)
		}
	}()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
