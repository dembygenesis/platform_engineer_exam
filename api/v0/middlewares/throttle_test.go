package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var m sync.RWMutex
var wg sync.WaitGroup

func TestValidate_Throttle(t *testing.T) {
	app := fiber.New()
	app.Get("/:token/validate", Throttle())
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

					require.Equal(t, helpers.WrapErrInErrMap(ErrThrottleLimitExceeded), errorRespExpected)
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
