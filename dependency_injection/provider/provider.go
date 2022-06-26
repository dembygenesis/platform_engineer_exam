package provider

import (
	BusinessToken "github.com/dembygenesis/platform_engineer_exam/business/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	PersistenceMYSQL "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	PersistenceToken "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/user"
	"github.com/pkg/errors"
	"github.com/sarulabs/dingo/v4"
	"log"
)

type Provider struct {
	dingo.BaseProvider
}

const (
	configLayer                = "config"
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
				cfg, err := config.NewConfig(".env")
				if err != nil {
					log.Fatalf("error setting up the config layer: :%v", err.Error())
				}
				return cfg, nil
			},
		},
		{
			Name: mysqlConnection,
			Build: func(config *config.Config) (*PersistenceMYSQL.MYSQLConnection, error) {
				mysql, err := PersistenceMYSQL.NewMYSQLConnection(config.DatabaseCredentials)
				if err != nil {
					log.Fatalf("error establishing connection to MYSQL: :%v", err.Error())
				}
				return mysql, err
			},
		},
		{
			Name: mysqlTokenPersistenceLayer,
			Build: func(connection *PersistenceMYSQL.MYSQLConnection) (*PersistenceToken.PersistenceToken, error) {
				return PersistenceToken.NewPersistenceToken(connection.DB), nil
			},
		},
		{
			Name: mysqlUserPersistenceLayer,
			Build: func(connection *PersistenceMYSQL.MYSQLConnection) (*user.PersistenceUser, error) {
				return user.NewPersistenceUser(connection.DB), nil
			},
		},
		{
			Name: businessToken,
			Build: func(config *config.Config, persistenceToken *PersistenceToken.PersistenceToken) (*BusinessToken.BusinessToken, error) {
				return BusinessToken.NewBusinessToken(persistenceToken, config.App.TokenDaysValid), nil
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
