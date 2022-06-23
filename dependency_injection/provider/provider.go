package provider

import (
	"fmt"
	BusinessToken "github.com/dembygenesis/platform_engineer_exam/business/token"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/token"
	"github.com/pkg/errors"
	"github.com/sarulabs/dingo/v4"
	"log"
)

type Provider struct {
	dingo.BaseProvider
}

const (
	configLayer                = "config"
	mysqlTokenPersistenceLayer = "mysql_token_persistence"
	businessToken              = "business_token"
)

// getServices is the main configuration func that produces the singleton
func getServices() (*[]dingo.Def, error) {
	var Services = []dingo.Def{
		{
			Name: configLayer,
			Build: func() (*config.Config, error) {
				return config.NewConfig(".env")
			},
		},
		{
			Name: mysqlTokenPersistenceLayer,
			Build: func(config *config.Config) (*token.MYSQLPersistence, error) {
				persistence, err := token.NewMYSQLPersistence(
					config.DatabaseCredentials.Host,
					config.DatabaseCredentials.User,
					config.DatabaseCredentials.Pass,
					config.DatabaseCredentials.Database,
					config.DatabaseCredentials.Port,
				)
				if err != nil {
					log.Fatalf("error establishing the mysql persistence: %v", err.Error())
				}
				fmt.Println("persistence", persistence)
				return persistence, nil
			},
		},
		{
			Name: businessToken,
			Build: func(mysqlPersistence *token.MYSQLPersistence) (*BusinessToken.BusinessToken, error) {
				return BusinessToken.NewBusinessToken(mysqlPersistence), nil
			},
		},
	}
	return &Services, nil
}

// Load bootstrap the dependencies
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
