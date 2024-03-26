package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"weather-station/internal/client"
	"weather-station/internal/metrics"
)

type CurrentWeather struct {
	Coord       Coord     `json:"coord"`
	Weather     []Weather `json:"weather"`
	Base        string    `json:"base"`
	Main        Main      `json:"main"`
	Visibility  int       `json:"visibility"`
	Wind        Wind      `json:"wind"`
	Rain        Rain      `json:"rain"`
	Clouds      Clouds    `json:"clouds"`
	Dt          int       `json:"dt"`
	Sys         Sys       `json:"sys"`
	Timezone    int       `json:"timezone"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Cod         int       `json:"cod"`
	ZipCode     string
	Logger      *slog.Logger
	client      client.HTTPGetter
	metrics     metrics.Observable
	subscribers []Subscriber
}
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}
type Rain struct {
	OneH float64 `json:"1h"`
}
type Clouds struct {
	All int `json:"all"`
}
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

func New(httpClient client.HTTPGetter, zipCode string, logger *slog.Logger, metrics metrics.Observable) *CurrentWeather {
	return &CurrentWeather{
		ZipCode: zipCode,
		Logger:  logger,
		client:  httpClient,
		metrics: metrics,
	}
}

func (w *CurrentWeather) GetTemperature() float64 {
	return w.Main.Temp
}

type Subscriber interface {
	Notify(weather *CurrentWeather) error
}

func (w *CurrentWeather) AddSubscriber(s Subscriber) {
	w.subscribers = append(w.subscribers, s)
}

func (w *CurrentWeather) Call() {
	if err := w.Get(); err != nil {
		w.Logger.Error("Hit some problem", "Error", err)
	}

	for _, subscriber := range w.subscribers {
		if err := subscriber.Notify(w); err != nil {
			w.Logger.Error("Unable to notify", "Error", err)
		}
	}
}

func (w *CurrentWeather) Get() error {
	w.Logger.Info("made http call")
	response, err := w.client.Get(w.ZipCode)
	if err != nil {
		return fmt.Errorf("HTTP client unable to make call, %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		w.metrics.ObserveAPIError(fmt.Sprint(response.StatusCode))
		return handleBadResponse(response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&w)

	w.Logger.Info("values", "weather", w)

	if err != nil {
		return fmt.Errorf("error decoding the response body, %v", err)
	}

	w.metrics.ObserveSuccess(w.Main.Temp, w.ZipCode, w.Name)

	return nil
}

func handleBadResponse(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return errors.New("bad API key")
	case http.StatusTooManyRequests:
		return errors.New("rate limited")
	}
	return fmt.Errorf("unknown api error, status code: %d", statusCode)
}
