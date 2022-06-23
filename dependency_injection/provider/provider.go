package provider

import (
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/pkg/errors"
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

const (
	configLayer = "config"
)

func getServices() (*[]dingo.Def, error) {
	var Services = []dingo.Def{
		{
			Name: configLayer,
			Build: func() (*config.Config, error) {
				return config.NewConfig()
			},
		},
	}
	return &Services, nil
}

func (p *Provider) Load() error {
	services, err := getServices()
	if err != nil {
		return errors.Wrap(err, "error trying to load the provider")
	}

	err = p.AddDefSlice(*services)
	if err != nil {
		return errors.Wrap(err, "error adding dependency definitions")
	}
	return nil
}
