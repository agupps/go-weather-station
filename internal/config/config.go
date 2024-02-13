package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Name              string            `yaml:"name"`
	Secret            string            `yaml:"secret"`
	WeatherProperties WeatherProperties `yaml:"weatherProperties"`
	Locations         []string          `yaml:"locations"`
}

type WeatherProperties struct {
	Units    string `yaml:"units"`
	Language string `yaml:"language"`
}

func (c *Config) Parse() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
