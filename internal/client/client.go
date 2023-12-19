package client

import (
	"fmt"
	"net/http"
	"os"
)

type HTTPGetter interface {
	Get(string) (*http.Response, error)
}

type OpenWeatherMapClient struct {
	apiKey string
}

func NewOpenWeatherMapClient() *OpenWeatherMapClient {
	return &OpenWeatherMapClient{
		apiKey: apiKey,
	}
}

var apiKey = os.Getenv("OWM_API_KEY")

var baseURL = "https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s"


func (client OpenWeatherMapClient) Get(zipCode string) (*http.Response, error) {
	url := fmt.Sprintf(baseURL, zipCode, client.apiKey)
	return http.Get(url)
}
