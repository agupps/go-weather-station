package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

type Locations struct {
	dao *daos.Dao
}

type Location struct {
	Zipcode     string  `db:zipcode`
	Name        string  `db:name`
	Temperature float64 `db:temperature`
}

func NewLocations(dao *daos.Dao) *Locations {
	return &Locations{dao: dao}
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

func (l *Locations) Update(location Location) error {
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
func (l *Locations) Delete(location Location) error {
	if _, err := l.dao.DB().Delete(
		"locations",
		&dbx.HashExp{"zipcode": location.Zipcode},
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
