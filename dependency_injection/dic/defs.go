package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	token1 "github.com/dembygenesis/platform_engineer_exam/business/token"
	config "github.com/dembygenesis/platform_engineer_exam/src/config"
	token "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/token"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "business_token",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("business_token")
				if err != nil {
					var eo *token1.BusinessToken
					return eo, err
				}
				pi0, err := ctn.SafeGet("mysql_token_persistence")
				if err != nil {
					var eo *token1.BusinessToken
					return eo, err
				}
				p0, ok := pi0.(*token.MYSQLPersistence)
				if !ok {
					var eo *token1.BusinessToken
					return eo, errors.New("could not cast parameter 0 to *token.MYSQLPersistence")
				}
				b, ok := d.Build.(func(*token.MYSQLPersistence) (*token1.BusinessToken, error))
				if !ok {
					var eo *token1.BusinessToken
					return eo, errors.New("could not cast build function to func(*token.MYSQLPersistence) (*token1.BusinessToken, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
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
		{
			Name:  "mysql_token_persistence",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("mysql_token_persistence")
				if err != nil {
					var eo *token.MYSQLPersistence
					return eo, err
				}
				pi0, err := ctn.SafeGet("config")
				if err != nil {
					var eo *token.MYSQLPersistence
					return eo, err
				}
				p0, ok := pi0.(*config.Config)
				if !ok {
					var eo *token.MYSQLPersistence
					return eo, errors.New("could not cast parameter 0 to *config.Config")
				}
				b, ok := d.Build.(func(*config.Config) (*token.MYSQLPersistence, error))
				if !ok {
					var eo *token.MYSQLPersistence
					return eo, errors.New("could not cast build function to func(*config.Config) (*token.MYSQLPersistence, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
	}
}
