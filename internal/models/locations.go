package models

import (
	"fmt"
	"net/http"
)

type Location struct {
	Zipcode     string  `json:"zipcode"`
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
}

func NewLocation(zipcode string) *Location {
	return &Location{
		Zipcode:     zipcode,
		Name:        "",
		Temperature: float64(0),
	}
}

func (l *Location) Get() error {
	l.Logger.Info("made http call")
	response, err := l.client.Get(l.ZipCode)
	if err != nil {
		return fmt.Errorf("HTTP client unable to make call, %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		l.metrics.ObserveAPIError(fmt.Sprint(response.StatusCode))
		return handleBadResponse(response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&l)

	l.Logger.Info("values", "weather", l)

	if err != nil {
		return fmt.Errorf("error decoding the response body, %v", err)
	}

	l.metrics.ObserveSuccess(l.Main.Temp, l.ZipCode, l.Name)

	return nil
}
