package weather

import (
	"errors"
	"fmt"
	"http"
	"json"
	"os"
)

type Weather struct {
	zipCode string
	countryCode string
	GeoPos   Coordinates `json:"coord"`
	Sys      Sys         `json:"sys"`
	Base     string      `json:"base"`
	Weather  []WeatherMain   `json:"weather"`
	Main     Main        `json:"main"`
	Visibility int       `json:"visibility"`
	Wind     Wind        `json:"wind"`
	Clouds   Clouds      `json:"clouds"`
	Rain     Rain        `json:"rain"`
	Snow     Snow        `json:"snow"`
	Dt       int         `json:"dt"`
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Cod      int         `json:"cod"`
	Timezone int         `json:"timezone"`
	Unit     string
	Lang     string
	Key      string
	*Settings
}

// Coordinates struct holds longitude and latitude data in returned
// JSON or as parameter data for requests using longitude and latitude.
type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

// Sys struct contains general information about the request
// and the surrounding area for where the request was made.
type Sys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

// Wind struct contains the speed and degree of the wind.
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

// Weather struct holds high-level, basic info on the returned
// data.
type WeatherMain struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Clouds struct holds data regarding cloud cover.
type Clouds struct {
	All int `json:"all"`
}

// Rain struct contains 3 hour data
type Rain struct {
	OneH   float64 `json:"1h,omitempty"`
	ThreeH float64 `json:"3h,omitempty"`
}

// Snow struct contains 3 hour data
type Snow struct {
	OneH   float64 `json:"1h,omitempty"`
	ThreeH float64 `json:"3h,omitempty"`
}

// Main struct contains the temperates, humidity, pressure for the request.
type Main struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
}

// Settings holds the client settings
type Settings struct {
	client *http.Client
}

var apiKey = os.Getenv("OWM_API_KEY")

func New(zipCode string) *Weather {
	return &Weather{
		zipCode: zipCode,
		Key: apiKey,
		countryCode: "US",
		Unit: "F",
		Lang: "EN",
	}
}

func (w *Weather) getTemperature() float64 {
	return w.Main.Temp
}

var (
	baseURL = "https://api.openweathermap.org/data/2.5/weather?%s"
	errInvalidKey          = errors.New("invalid api key")
)


func (w *Weather) GetWeather() {
	response, err := w.client.Get(fmt.Sprintf(fmt.Sprintf(baseURL, "appid=%s&zip=%s,%s&units=%s&lang=%s"), w.Key, w.zipCode, w.countryCode, w.Unit, w.Lang))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return errInvalidKey
	}

	w = json.NewDecoder(response.Body).Decode(&w)
}