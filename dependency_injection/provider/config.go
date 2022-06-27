package provider

import (
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/sarulabs/dingo/v4"
	"log"
)

const (
	configLayer = "config"
)

func getConfigLayers() *[]dingo.Def {
	return &[]dingo.Def{
		{
			Name: configLayer,
			Build: func() (*config.Config, error) {
				cfg, err := config.NewConfig(".env")
				if err != nil {
					log.Fatalf("error setting up the config layer: :%v", err.Error())
				}
				return cfg, nil
			},
		},
	}
}
