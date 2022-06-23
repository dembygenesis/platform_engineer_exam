package middlewares

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ProtectedRoute guards a route using the "Basic Auth" protocol
func ProtectedRoute(ctn *dic.Container) func(c *fiber.Ctx) error {
	userPersistence := ctn.GetMysqlUserPersistence()
	logger := ctn.GetLogger()

	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			matched, err := userPersistence.BasicAuth(user, pass)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("error_basic_auth")
				return false
			}
			if !matched {
				logger.WithFields(logrus.Fields{
					"user": user,
				}).Info("info_no_match")
				return false
			}
			return true
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(http.StatusUnauthorized).JSON(helpers.WrapStrInErrMap("Unauthorized"))
		},
	})
}
