package token

import "github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"

type dataPersistence interface {
	// Generate creates a new 6-12 digit authentication token
	Generate(createdBy int) (string, error)

	// Validate checks if a string is registered
	Validate(s string) error
}

type BusinessToken struct {
	dataLayer dataPersistence
}

func (b *BusinessToken) Generate(user *models_schema.User) (string, error) {
	return b.dataLayer.Generate(user.ID)
}

func (b *BusinessToken) Validate(s string) error {
	return b.dataLayer.Validate(s)
}

func NewBusinessToken(mysqlDataPersistence dataPersistence) *BusinessToken {
	return &BusinessToken{
		dataLayer: mysqlDataPersistence,
	}
}
