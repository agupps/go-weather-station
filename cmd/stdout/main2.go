package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
	"weather-station/internal/client"
	"weather-station/internal/config"
	"weather-station/internal/metrics"
	"weather-station/internal/poller"
	"weather-station/internal/weather"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	p := poller.NewPoller(time.Second, logger)
	c := &config.Config{}
	if err := c.Parse(); err != nil {
		panic(err)
	}

	registry := prometheus.NewRegistry()

	tempMetrics := metrics.NewMetrics(registry)

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	client := client.NewOpenWeatherMapClient(c.Secret, c.WeatherProperties)

	p.Add(weather.New(client, "21163", logger, tempMetrics))
	p.Add(weather.New(client, "20008", logger, tempMetrics))
	p.Add(weather.New(client, "27520", logger, tempMetrics))
	p.Add(weather.New(client, "95134", logger, tempMetrics))
	go p.Start()

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	log.Fatal(http.ListenAndServe(":8080", nil))

	// woodstock := weather.New("21163")
	// fmt.Printf("Weather response: %+v\n", woodstock.Get())
}
