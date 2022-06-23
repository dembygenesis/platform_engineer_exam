package auth

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"net/http"
)

// ProtectedRoute guards a route using the "Basic Auth" protocol
func ProtectedRoute(ctn *dic.Container) func(c *fiber.Ctx) error {

	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			userPersistence, err := ctn.SafeGetMysqlUserPersistence()
			if err != nil {
				fmt.Println("to get user persistence - convert to log, later")
				return false
			}
			matched, err := userPersistence.BasicAuth(user, pass)
			if err != nil {
				fmt.Println("to get user persistence - convert to log, later", err.Error())
				return false
			}
			if !matched {
				fmt.Println("No match!")
				return false
			}
			return true
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(http.StatusUnauthorized).JSON(helpers.WrapInErrMap("Unauthorized"))
		},
	})
}
