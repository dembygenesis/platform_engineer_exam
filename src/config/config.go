package config

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/validation"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

var (
	errValidatingStructParams         = errors.New("error validating struct params")
	errReturnedFromParamValidation    = errors.New("errors returned from param validation")
	errUnidentifiedTokenLapseDuration = errors.New("error - unidentified token lapse duration")
)

// DatabaseCredentials holds our database env settings
type DatabaseCredentials struct {
	Host     string `mapstructure:"DB_HOST" validate:"required"`
	User     string `mapstructure:"DB_USERNAME" validate:"required"`
	Pass     string `mapstructure:"DB_PASSWORD" validate:"required"`
	Port     string `mapstructure:"DB_PORT" validate:"required"`
	Database string `mapstructure:"DB_DATABASE" validate:"required"`
}

type App struct {
	TokenLapseSettings string `mapstructure:"TOKEN_LAPSE_DURATION" validate:"required"`
	TokenLapseDuration float64
}

type API struct {
	Port int `mapstructure:"PORT"`
}

type Config struct {
	DatabaseCredentials DatabaseCredentials
	API                 API
	App                 App
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

	err = viper.Unmarshal(&config.App)
	if err != nil {
		return config, errors.Wrap(err, "error unmarshalling the app")
	}

	configStructs := []interface{}{
		config.DatabaseCredentials,
		config.API,
		config.App,
	}
	for _, configStruct := range configStructs {
		errs, err := validation.ValidateStructParams(&configStruct)
		if err != nil {
			return nil, errors.Wrap(err, errValidatingStructParams.Error())
		}
		if len(errs) > 0 {
			return nil, errors.Wrap(errors.New(strings.Join(errs[:], ",")),
				errReturnedFromParamValidation.Error())
		}
	}

	switch config.App.TokenLapseSettings {
	case "7 Days":
		fmt.Println("Yo 7 days")
	case "30 Seconds":
		fmt.Println("Yo 30 Seconds")
	default:
		return nil, errUnidentifiedTokenLapseDuration
	}

	return config, nil
}
