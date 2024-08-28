package models

import (
	"weather-station/internal/weather"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

type Locations struct {
	dao *daos.Dao
}

type Location struct {
	Zipcode     string  `json:"zipcode"`
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
}

func (l *Locations) Notify(w *weather.CurrentWeather) error {
	if err := l.Update(fromCurrentWeather(w)); err != nil {
		return err
	}
	return nil
}

func fromCurrentWeather(w *weather.CurrentWeather) *Location {
	return &Location{Zipcode: w.ZipCode, Name: w.Name, Temperature: w.GetTemperature()}
}

func NewLocations(dao *daos.Dao) *Locations {
	return &Locations{dao: dao}
}

func NewLocation(zipcode string) *Location {
	return &Location{
		Zipcode:     zipcode,
		Name:        "",
		Temperature: float64(0),
	}
}

func (l *Locations) Create(location *Location) error {
	if _, err := l.dao.DB().Insert("locations",
		dbx.Params{
			"zipcode":     location.Zipcode,
			"temperature": location.Temperature,
			"name":        location.Name,
		}).Execute(); err != nil {
		return err
	}
	return nil
}

func (l *Locations) Update(location *Location) error {
	if _, err := l.dao.DB().Update(
		"locations",
		dbx.Params{
			"temperature": location.Temperature,
			"name":        location.Name,
		},
		&dbx.HashExp{"zipcode": location.Zipcode},
	).Execute(); err != nil {
		return err
	}
	return nil
}
func (l *Locations) Delete() error {
	if _, err := l.dao.DB().Delete(
		"locations",
		nil,
	).Execute(); err != nil {
		return err
	}
	return nil

}
func (l *Locations) All() ([]Location, error) {
	allLocations := &[]Location{}
	err := l.dao.DB().NewQuery("SELECT * FROM locations").All(allLocations)

	return *allLocations, err
}
