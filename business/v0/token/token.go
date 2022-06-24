package token

import (
	"context"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
)

type dataPersistence interface {
	GetAll(ctx context.Context) ([]models_schema.Token, error)

	// Generate creates a new 6-12 digit authentication on token
	Generate(ctx context.Context, createdBy int) (string, error)

	// Validate checks if a string is registered
	Validate(ctx context.Context, str string, lapseLimit float64) error
}

type BusinessToken struct {
	dataLayer dataPersistence
}

func (b *BusinessToken) GetAll(ctx context.Context) ([]models_schema.Token, error) {
	return b.dataLayer.GetAll(ctx)
}

func (b *BusinessToken) Generate(ctx context.Context, user *models_schema.User) (string, error) {
	return b.dataLayer.Generate(ctx, user.ID)
}

func (b *BusinessToken) Validate(ctx context.Context, str string, lapseLimit float64) error {
	return b.dataLayer.Validate(ctx, str, lapseLimit)
}

func NewBusinessToken(mysqlDataPersistence dataPersistence) *BusinessToken {
	return &BusinessToken{
		dataLayer: mysqlDataPersistence,
	}
}
