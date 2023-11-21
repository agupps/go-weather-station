package main

import (
	"log/slog"
	"os"
	"time"
	"weather-station/poller"
	"weather-station/weather"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))


	p := poller.NewPoller(time.Second, logger)
	p.Add(weather.New("21163", logger))
	p.Add(weather.New("20008", logger))
	p.Add(weather.New("27520", logger))
	p.Add(weather.New("95134", logger))
	p.Start()




	// woodstock := weather.New("21163")
	// fmt.Printf("Weather response: %+v\n", woodstock.Get())
}
