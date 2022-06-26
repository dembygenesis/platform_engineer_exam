package token

import (
	"context"
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/business/v0/token/tokenfakes"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/stretchr/testify/assert"
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

func TestBusinessToken_Validate_FailPath_UpdateTokenToExpired(t *testing.T) {
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
	}, errUpdateTokenToExpired)

	businessToken := NewBusinessToken(&fakeDataPersistence, 7)
	err := businessToken.Validate(context.Background(), tokenKey, 0, "")
	t.Run("Test Validate - Update Token To Expired", func(t *testing.T) {
		defer func() {
			fmt.Println("=============err=============", err)
			require.Error(t, err)
		}()
	})
}

func TestBusinessToken_Validate_FailPath_GetToken(t *testing.T) {
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
	}, errGetToken)

	businessToken := NewBusinessToken(&fakeDataPersistence, 7)
	err := businessToken.Validate(context.Background(), tokenKey, 0, "")
	t.Run("Test Validate - Happy Path", func(t *testing.T) {
		require.NoError(t, err)

		errMsg := err.Error()
		wantErrMsg := errGetToken.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func TestBusinessToken_Validate_FailPath_Revoked(t *testing.T) {
	tokenKey := "123456"

	fakeDataPersistence := tokenfakes.FakeDataPersistence{}
	fakeDataPersistence.GetTokenReturns(&models.Token{
		Id:        1,
		Key:       tokenKey,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Revoked:   true,
		Expired:   false,
		CreatedBy: "Demby",
	}, nil)

	businessToken := NewBusinessToken(&fakeDataPersistence, 7)
	err := businessToken.Validate(context.Background(), tokenKey, 0, "")
	t.Run("Test Validate - Fail Path Revoked", func(t *testing.T) {
		require.NoError(t, err)

		errMsg := err.Error()
		wantErrMsg := errTokenRevoked.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func TestBusinessToken_Validate_FailPath_Expired(t *testing.T) {
	tokenKey := "123456"

	fakeDataPersistence := tokenfakes.FakeDataPersistence{}
	fakeDataPersistence.GetTokenReturns(&models.Token{
		Id:        1,
		Key:       tokenKey,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Revoked:   false,
		Expired:   true,
		CreatedBy: "Demby",
	}, nil)

	businessToken := NewBusinessToken(&fakeDataPersistence, 7)
	err := businessToken.Validate(context.Background(), tokenKey, 0, "")
	t.Run("Test Validate - Fail Path Expired", func(t *testing.T) {
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errTokenExpired.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func TestBusinessToken_Validate_FailPath_DeterminedExpired(t *testing.T) {
	tokenKey := "123456"

	fakeDataPersistence := tokenfakes.FakeDataPersistence{}
	fakeDataPersistence.GetTokenReturns(&models.Token{
		Id:        1,
		Key:       tokenKey,
		CreatedAt: time.Now().AddDate(0, 0, -7),
		ExpiresAt: time.Now(),
		Revoked:   false,
		Expired:   false,
		CreatedBy: "Demby",
	}, nil)

	businessToken := NewBusinessToken(&fakeDataPersistence, 7)
	err := businessToken.Validate(context.Background(), tokenKey, 0, "")
	fmt.Println("err err err", err)
	t.Run("Test Validate - Fail Path Determined Expired", func(t *testing.T) {
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errTokenDeterminedExpired.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}
