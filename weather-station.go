package main

import (
	"fmt"
	"log"
	"os"

	// Shortening the import reference name seems to make it a bit easier
	owm "github.com/briandowns/openweathermap"
)

var apiKey = os.Getenv("OWM_API_KEY")

func main() {
	w, err := owm.NewCurrent("F", "ru", apiKey) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName("Phoenix")
	fmt.Println(w.GeoPos.Latitude)
	fmt.Println(w.GeoPos.Longitude)
	fmt.Println(w.Sys.Country)
	fmt.Println(w.Main.Temp)
}
