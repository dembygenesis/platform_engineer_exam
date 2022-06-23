package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// DatabaseCredentials holds our database env settings
type DatabaseCredentials struct {
	User    string `mapstructure:"DB_USERNAME"`
	Pass    string `mapstructure:"DB_PASSWORD"`
	Port    string `mapstructure:"DB_PORT"`
	Schema  string `mapstructure:"DB_SCHEMA"`
	Timeout string `mapstructure:"DB_TIMEOUT"`
}

type API struct {
	Port int `mapstructure:"PORT"`
}

type Config struct {
	DatabaseCredentials DatabaseCredentials
	API                 API
}

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error reading config file")
	}
	config := &Config{}

	// Database
	err = viper.Unmarshal(&config.DatabaseCredentials)
	if err != nil {
		return config, errors.Wrap(err, "error trying to unmarshal the database credentials")
	}

	// Port
	err = viper.Unmarshal(&config.API)
	if err != nil {
		return config, errors.Wrap(err, "error unmarshalling the port")
	}

	return config, nil
}
