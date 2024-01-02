package client

import (
	"fmt"
	"net/http"
	"weather-station/internal/config"
)

type HTTPGetter interface {
	Get(string) (*http.Response, error)
}

type OpenWeatherMapClient struct {
	apiKey string
	lang   string
	units  string
}

func NewOpenWeatherMapClient(apiKey string, weatherProperties config.WeatherProperties) *OpenWeatherMapClient {
	return &OpenWeatherMapClient{
		apiKey: apiKey,
		units:  weatherProperties.Units,
		lang:   weatherProperties.Language,
	}
}

var baseURL = "https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s&units=%s&lang=%s"

func (client OpenWeatherMapClient) Get(zipCode string) (*http.Response, error) {
	url := fmt.Sprintf(baseURL, zipCode, client.apiKey, client.units, client.lang)
	return http.Get(url)
}
