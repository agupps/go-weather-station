package main

import (
	"os"
	"weather-station/internal/app"
)

func main() {
	os.Exit(app.New().Run())
}
