package token

import (
	"context"
	"github.com/dembygenesis/platform_engineer_exam/business/v0/token/tokenfakes"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBusinessToken_Validate_HappyPath(t *testing.T) {
	tokenKey := "123456"

	fakeDataPersistence := tokenfakes.FakeDataPersistence{}
	fakeDataPersistence.GetTokenReturns(&models.Token{
		Id:        1,
		Key:       tokenKey,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Revoked:   false,
		Expired:   false,
		CreatedBy: "Demby",
	}, nil)

	businessToken := NewBusinessToken(&fakeDataPersistence, 7)
	err := businessToken.Validate(context.Background(), tokenKey, 0, "")
	t.Run("Test Validate - Happy Path", func(t *testing.T) {
		require.NoError(t, err)
	})
}
