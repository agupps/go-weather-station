package main

import (
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	owm "github.com/briandowns/openweathermap"
)

func main() {
	a := app.New()
	w := a.NewWindow("Go Weather Station")

	entry := widget.NewEntry()
	outputTemperature := widget.NewLabel("")
	outputSunsetTime := widget.NewLabel("")

	topLabel := widget.NewLabel("Welcome! Please input a city to get the current temperature!")
	w.SetContent(container.NewVBox(
		topLabel,
		entry,
		widget.NewButton("Submit", func() {
			outputTemperature.SetText("Temperature: " + strconv.FormatFloat(getTemperatureByLocation(entry.Text), 'f', -1, 64) + " F")
			outputSunsetTime.SetText("Humidity: " + strconv.FormatInt(int64(getHumidityByLocation(entry.Text)), 10) + " %")
		}),
		outputTemperature,
		outputSunsetTime,
	))

	w.ShowAndRun()
}

var apiKey = os.Getenv("OWM_API_KEY")

func getTemperatureByLocation(location string) float64 {
	w, err := owm.NewCurrent("F", "EN", apiKey) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName(location)
	return w.Main.Temp
}

func getHumidityByLocation(location string) int {
	w, err := owm.NewCurrent("F", "EN", apiKey) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName(location)
	return w.Main.Humidity
}
