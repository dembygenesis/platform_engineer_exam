package token

import (
	"context"
	"github.com/dembygenesis/platform_engineer_exam/api/v0/token/tokenfakes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidate_StatusOk(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.ValidateReturns(nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/:token/validate", apiToken.ValidateToken)

	req := httptest.NewRequest("GET", "/mock_token_value/validate", nil)

	ctx := context.Background()
	req.WithContext(context.WithValue(ctx, "token", "12345"))

	resp, _ := app.Test(req, 1)
	t.Run("Test Validate - StatusOk", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestValidate_InternalServerError(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.ValidateReturns(errMockValidate)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/:token/validate", apiToken.ValidateToken)

	req := httptest.NewRequest("GET", "/mock_token_value/validate", nil)

	ctx := context.Background()
	req.WithContext(context.WithValue(ctx, "token", "12345"))

	resp, _ := app.Test(req, 1)
	t.Run("Test Validate - Internal Server Error", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestRevoke_StatusOk(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.RevokeReturns(nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/:token/revoke", apiToken.Revoke)

	req := httptest.NewRequest("GET", "/mock_token_value/revoke", nil)

	ctx := context.Background()
	req.WithContext(context.WithValue(ctx, "token", "12345"))

	resp, _ := app.Test(req, 1)
	t.Run("Test GetAll - StatusOk", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestRevoke_InternalServerError(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.RevokeReturns(errMockRevoke)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/:token/revoke", apiToken.Revoke)

	req := httptest.NewRequest("GET", "/mock_token_value/revoke", nil)

	ctx := context.Background()
	req.WithContext(context.WithValue(ctx, "token", "12345"))

	resp, _ := app.Test(req, 1)
	t.Run("Test GetAll - Internal Server Error", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestGetAll_StatusOk(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.GetAllReturns(nil, nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/", apiToken.GetAll)

	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test GetAll - Ok", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGetAll_InternalServerError(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.GetAllReturns(nil, errMockGetAll)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/", apiToken.GetAll)

	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test GetAll - Internal Server Error", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
