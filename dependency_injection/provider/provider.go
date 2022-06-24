package provider

import (
	BusinessToken "github.com/dembygenesis/platform_engineer_exam/business/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	v0 "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	token2 "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/user"
	"github.com/pkg/errors"
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

const (
	configLayer = "config"
	// logger                     = "logger"
	mysqlConnection            = "mysql_connection"
	mysqlTokenPersistenceLayer = "mysql_token_persistence"
	mysqlUserPersistenceLayer  = "mysql_user_persistence"
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
			Name: mysqlConnection,
			Build: func(config *config.Config) (*v0.MYSQLConnection, error) {
				return v0.NewMYSQLConnection(config.DatabaseCredentials)
			},
		},
		{
			Name: mysqlTokenPersistenceLayer,
			Build: func(connection *v0.MYSQLConnection) (*token2.PersistenceToken, error) {
				return token2.NewPersistenceToken(connection.DB), nil
			},
		},
		{
			Name: mysqlUserPersistenceLayer,
			Build: func(connection *v0.MYSQLConnection) (*user.PersistenceUser, error) {
				return user.NewPersistenceUser(connection.DB), nil
			},
		},
		{
			Name: businessToken,
			Build: func(persistenceToken *token2.PersistenceToken) (*BusinessToken.BusinessToken, error) {
				return BusinessToken.NewBusinessToken(persistenceToken), nil
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
