package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// DatabaseCredentials holds our database env settings
type DatabaseCredentials struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USERNAME"`
	Pass     string `mapstructure:"DB_PASSWORD"`
	Port     string `mapstructure:"DB_PORT"`
	Database string `mapstructure:"DB_DATABASE"`
}

type API struct {
	Port int `mapstructure:"PORT"`
}

type Config struct {
	DatabaseCredentials DatabaseCredentials
	API                 API
}

// NewConfig reads values from the .env file, and writes them to the Config struct
func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error reading config file")
	}
	config := &Config{}

	err = viper.Unmarshal(&config.DatabaseCredentials)
	if err != nil {
		return config, errors.Wrap(err, "error trying to unmarshal the database credentials")
	}

	err = viper.Unmarshal(&config.API)
	if err != nil {
		return config, errors.Wrap(err, "error unmarshalling the port")
	}

	return config, nil
}
