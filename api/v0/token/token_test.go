package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/v0/token/tokenfakes"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAll_StatusOk(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.GetAllReturns(nil, nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/", apiToken.GetAll)

	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test GetAll - Ok", func(t *testing.T) {
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGetAll_InternalServerError(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.GetAllReturns(nil, errors.New("err getting fake biz function"))

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/", apiToken.GetAll)

	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test GetAll - Internal Server Error", func(t *testing.T) {
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
