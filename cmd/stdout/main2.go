package main

import (
	"log/slog"
	"os"
	"time"
	"weather-station/internal/client"
	"weather-station/internal/poller"
	"weather-station/internal/weather"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))


	p := poller.NewPoller(time.Second, logger)
	client := client.NewOpenWeatherMapClient()
	p.Add(weather.New(client, "21163", logger))
	p.Add(weather.New(client, "20008", logger))
	p.Add(weather.New(client, "27520", logger))
	p.Add(weather.New(client, "95134", logger))
	p.Start()




	// woodstock := weather.New("21163")
	// fmt.Printf("Weather response: %+v\n", woodstock.Get())
}
