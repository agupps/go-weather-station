package weather_test

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"weather-station/internal/weather"

	"github.com/stretchr/testify/assert"
)

type HTTPMock struct {
	response *http.Response
}

func (client HTTPMock) Get(_ string) (*http.Response, error) {
	return client.response, nil
}

type MetricsMock struct {
}

func (m MetricsMock) ObserveSuccess(_ float64, _ ...string) {

}

func (m MetricsMock) ObserveAPIError(_ string) {

}

var sampleResponse = `
{
	"coord": {
	  "lon": 10.99,
	  "lat": 44.34
	},
	"weather": [
	  {
		"id": 501,
		"main": "Rain",
		"description": "moderate rain",
		"icon": "10d"
	  }
	],
	"base": "stations",
	"main": {
	  "temp": 298.48,
	  "feels_like": 298.74,
	  "temp_min": 297.56,
	  "temp_max": 300.05,
	  "pressure": 1015,
	  "humidity": 64,
	  "sea_level": 1015,
	  "grnd_level": 933
	},
	"visibility": 10000,
	"wind": {
	  "speed": 0.62,
	  "deg": 349,
	  "gust": 1.18
	},
	"rain": {
	  "1h": 3.16
	},
	"clouds": {
	  "all": 100
	},
	"dt": 1661870592,
	"sys": {
	  "type": 2,
	  "id": 2075663,
	  "country": "IT",
	  "sunrise": 1661834187,
	  "sunset": 1661882248
	},
	"timezone": 7200,
	"id": 3163858,
	"name": "Zocca",
	"cod": 200
  }                        
`

var badAPIKeyResponse = `{"cod":401, "message": "Invalid API key. Please see https://openweathermap.org/faq#error401 for more info."}`

func TestCurrentWeather_Call(t *testing.T) {
	testCases := []struct {
		name          string
		statusCode    int
		responseText  string
		errorExpected bool
	}{
		{
			name:          "successful call",
			statusCode:    http.StatusOK,
			responseText:  sampleResponse,
			errorExpected: false,
		},
		{
			name:          "bad api key",
			statusCode:    http.StatusUnauthorized,
			responseText:  badAPIKeyResponse,
			errorExpected: true,
		},
	}
	for _, testCase := range testCases {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		stringReader := strings.NewReader(testCase.responseText)

		exampleResponse := &http.Response{
			Body:       io.NopCloser(stringReader),
			Status:     http.StatusText(testCase.statusCode),
			StatusCode: testCase.statusCode,
		}
		mockClient := HTTPMock{response: exampleResponse}

		metricsMock := &MetricsMock{}

		w := weather.New(mockClient, "21163", logger, metricsMock)
		err := w.Get()
		if testCase.errorExpected {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, 298.48, w.Main.Temp)
		}
	}
}
