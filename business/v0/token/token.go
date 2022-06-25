package token

import (
	"context"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"time"
)

type dataPersistence interface {
	GetAll(ctx context.Context) ([]models.Token, error)

	// Generate creates a new 6-12 digit authentication on token
	Generate(ctx context.Context, createdBy int, randomStringsOverride string, createdAtOverride *time.Time) (string, error)

	// Validate checks if a string is registered
	Validate(ctx context.Context, str string, lapseLimit float64, lapseType string) error
}
type BusinessToken struct {
	dataLayer dataPersistence
}

func (b *BusinessToken) GetAll(ctx context.Context) ([]models.Token, error) {
	return b.dataLayer.GetAll(ctx)
}

func (b *BusinessToken) Generate(ctx context.Context, user *models_schema.User) (string, error) {
	return b.dataLayer.Generate(ctx, user.ID, "", nil)
}

func (b *BusinessToken) Validate(ctx context.Context, str string, lapseLimit float64, lapseType string) error {
	return b.dataLayer.Validate(ctx, str, lapseLimit, lapseType)
}

func NewBusinessToken(mysqlDataPersistence dataPersistence) *BusinessToken {
	return &BusinessToken{
		dataLayer: mysqlDataPersistence,
	}
}
