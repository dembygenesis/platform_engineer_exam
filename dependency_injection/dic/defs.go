package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	token1 "github.com/dembygenesis/platform_engineer_exam/business/v0/token"
	config "github.com/dembygenesis/platform_engineer_exam/src/config"
	mysql "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	token "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/token"
	user "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/user"
	logrus "github.com/sirupsen/logrus"
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
				p0, ok := pi0.(*token.PersistenceToken)
				if !ok {
					var eo *token1.BusinessToken
					return eo, errors.New("could not cast parameter 0 to *token.PersistenceToken")
				}
				b, ok := d.Build.(func(*token.PersistenceToken) (*token1.BusinessToken, error))
				if !ok {
					var eo *token1.BusinessToken
					return eo, errors.New("could not cast build function to func(*token.PersistenceToken) (*token1.BusinessToken, error)")
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
			Name:  "logger",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logger")
				if err != nil {
					var eo *logrus.Entry
					return eo, err
				}
				b, ok := d.Build.(func() (*logrus.Entry, error))
				if !ok {
					var eo *logrus.Entry
					return eo, errors.New("could not cast build function to func() (*logrus.Entry, error)")
				}
				return b()
			},
			Unshared: false,
		},
		{
			Name:  "mysql_connection",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("mysql_connection")
				if err != nil {
					var eo *mysql.MYSQLConnection
					return eo, err
				}
				pi0, err := ctn.SafeGet("config")
				if err != nil {
					var eo *mysql.MYSQLConnection
					return eo, err
				}
				p0, ok := pi0.(*config.Config)
				if !ok {
					var eo *mysql.MYSQLConnection
					return eo, errors.New("could not cast parameter 0 to *config.Config")
				}
				b, ok := d.Build.(func(*config.Config) (*mysql.MYSQLConnection, error))
				if !ok {
					var eo *mysql.MYSQLConnection
					return eo, errors.New("could not cast build function to func(*config.Config) (*mysql.MYSQLConnection, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "mysql_token_persistence",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("mysql_token_persistence")
				if err != nil {
					var eo *token.PersistenceToken
					return eo, err
				}
				pi0, err := ctn.SafeGet("mysql_connection")
				if err != nil {
					var eo *token.PersistenceToken
					return eo, err
				}
				p0, ok := pi0.(*mysql.MYSQLConnection)
				if !ok {
					var eo *token.PersistenceToken
					return eo, errors.New("could not cast parameter 0 to *mysql.MYSQLConnection")
				}
				b, ok := d.Build.(func(*mysql.MYSQLConnection) (*token.PersistenceToken, error))
				if !ok {
					var eo *token.PersistenceToken
					return eo, errors.New("could not cast build function to func(*mysql.MYSQLConnection) (*token.PersistenceToken, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
		{
			Name:  "mysql_user_persistence",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("mysql_user_persistence")
				if err != nil {
					var eo *user.PersistenceUser
					return eo, err
				}
				pi0, err := ctn.SafeGet("mysql_connection")
				if err != nil {
					var eo *user.PersistenceUser
					return eo, err
				}
				p0, ok := pi0.(*mysql.MYSQLConnection)
				if !ok {
					var eo *user.PersistenceUser
					return eo, errors.New("could not cast parameter 0 to *mysql.MYSQLConnection")
				}
				b, ok := d.Build.(func(*mysql.MYSQLConnection) (*user.PersistenceUser, error))
				if !ok {
					var eo *user.PersistenceUser
					return eo, errors.New("could not cast build function to func(*mysql.MYSQLConnection) (*user.PersistenceUser, error)")
				}
				return b(p0)
			},
			Unshared: false,
		},
	}
}
