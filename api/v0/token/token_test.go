package token

import (
	"bytes"
	"encoding/json"
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/api/v0/middlewares"
	"github.com/dembygenesis/platform_engineer_exam/api/v0/token/tokenfakes"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetToken_StatusCreated(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.GenerateReturns("12345", nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Post("/", func(ctx *fiber.Ctx) error {
		ctx.Locals("userMeta", &models.User{Id: 1})
		return ctx.Next()
	}, apiToken.GetToken)

	req := httptest.NewRequest("POST", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test Validate - StatusCreated", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

func TestGetToken_InternalServerError_GenerateFails(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.GenerateReturns("", errMockGenerate)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Post("/", func(ctx *fiber.Ctx) error {
		ctx.Locals("userMeta", &models.User{Id: 1})
		return ctx.Next()
	}, apiToken.GetToken)

	req := httptest.NewRequest("POST", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test Validate - Internal Server Error", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestGetToken_InternalServerError_UserMetaFails(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/", apiToken.GetToken)

	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := app.Test(req, 1)
	t.Run("Test Validate - Internal Server Error", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

var m sync.RWMutex
var wg sync.WaitGroup

func TestValidate_Throttle(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.ValidateReturns(nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/:token/validate", middlewares.Throttle(), apiToken.ValidateToken)
	req := httptest.NewRequest("GET", "/mock_token_value/validate", nil)

	iterator := 0

	t.Run("Test Validate - Throttled Requests", func(t *testing.T) {
		errorResponseCount := 0
		reqCount := 20

		wg.Add(reqCount)

		for reqCount > 0 {
			iterator++

			go func() {
				defer wg.Done()

				resp, err := app.Test(req, 1)
				require.NoError(t, err)

				readBody, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				err = resp.Body.Close()
				require.NoError(t, err)

				resp.Body = ioutil.NopCloser(bytes.NewReader(readBody))
				respStr := string(readBody)

				if resp.StatusCode == http.StatusForbidden {
					var errorRespExpected map[string][]string
					err = json.Unmarshal([]byte(respStr), &errorRespExpected)
					require.NoError(t, err)

					require.Equal(t, helpers.WrapErrInErrMap(middlewares.ErrThrottleLimitExceeded), errorRespExpected)
					m.Lock()
					errorResponseCount = errorResponseCount + 1
					m.Unlock()
				}
			}()
			reqCount--
		}
		wg.Wait()

		require.Equal(t, true, errorResponseCount > 4)
	})
}

func TestValidate_StatusOk(t *testing.T) {
	fakeBizFunctions := &tokenfakes.FakeBizFunctions{}
	fakeBizFunctions.ValidateReturns(nil)

	apiToken := NewAPIToken(fakeBizFunctions)

	app := fiber.New()
	app.Get("/:token/validate", apiToken.ValidateToken)

	req := httptest.NewRequest("GET", "/mock_token_value/validate", nil)

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
	app.Delete("/:token/revoke", apiToken.Revoke)

	req := httptest.NewRequest("DELETE", "/mock_token_value/revoke", nil)

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
	app.Delete("/:token/revoke", apiToken.Revoke)

	req := httptest.NewRequest("DELETE", "/mock_token_value/revoke", nil)

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
