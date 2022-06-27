package provider

import (
	BusinessToken "github.com/dembygenesis/platform_engineer_exam/business/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	PersistenceToken "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/token"
	"github.com/sarulabs/dingo/v4"
)

const (
	businessToken = "business_token"
)

func getBusinessLayers() *[]dingo.Def {
	return &[]dingo.Def{
		{
			Name: businessToken,
			Build: func(config *config.Config, persistenceToken *PersistenceToken.PersistenceToken) (*BusinessToken.BusinessToken, error) {
				return BusinessToken.NewBusinessToken(
					persistenceToken,
					config.App.TokenDaysValid,
					config.App.RandomCharMinLength,
					config.App.RandomCharMaxLength,
				), nil
			},
		},
	}
}
