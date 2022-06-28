package provider

import (
	"github.com/dembygenesis/platform_engineer_exam/api/v0/middlewares"
	"github.com/dembygenesis/platform_engineer_exam/api/v0/token"
	BusinessToken "github.com/dembygenesis/platform_engineer_exam/business/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/v0/user"
	"github.com/sarulabs/dingo/v4"
)

const (
	apiToken       = "api_token"
	apiMiddlewares = "api_middlewares"
)

func getAPILayers() *[]dingo.Def {
	return &[]dingo.Def{
		{
			Name: apiToken,
			Build: func(businessToken *BusinessToken.BusinessToken) (*token.APIToken, error) {
				return token.NewAPIToken(businessToken), nil
			},
		},
		{
			Name: apiMiddlewares,
			Build: func(user *user.PersistenceUser) (*middlewares.AuthRoutes, error) {
				return middlewares.NewAuthRoutes(user), nil
			},
		},
	}
}
