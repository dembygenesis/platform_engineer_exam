package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	config "github.com/dembygenesis/platform_engineer_exam/src/config"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "config",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("config")
				if err != nil {
					var eo *config.Config
					return eo, err
				}
				b, ok := d.Build.(func() (*config.Config, error))
				if !ok {
					var eo *config.Config
					return eo, errors.New("could not cast build function to func() (*config.Config, error)")
				}
				return b()
			},
			Unshared: false,
		},
	}
}
