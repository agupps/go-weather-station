package main

import (
	"fmt"
	"weather-station/weather"
)

func main() {
	w := weather.New("20008")
	fmt.Printf("Weather response: %+v\n", w.Get())

	woodstock := weather.New("21163")
	fmt.Printf("Weather response: %+v\n", woodstock.Get())
}
