package weather

import (
	"log"
	"os"

	owm "github.com/briandowns/openweathermap"
)

type Weather struct {
	client  *owm.CurrentWeatherData
	zipCode string
}

var apiKey = os.Getenv("OWM_API_KEY")

func New(zipCode string) Weather {
	w, err := owm.NewCurrent("F", "EN", apiKey)
	if err != nil {
		log.Fatalln(err)
	}
	return Weather{client: w, zipCode: zipCode}
}

func (w Weather) getTemperature() float64 {
	w.client.CurrentByZipcode(w.zipCode, "US")
	return w.client.Main.Temp
}

func (w Weather) Get() *owm.CurrentWeatherData {
	w.client.CurrentByZipcode(w.zipCode, "US")
	return w.client
}
