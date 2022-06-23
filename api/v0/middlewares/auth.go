package middlewares

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	user = "user"
	pass = "pass"
)

// ProtectedRoute guards a route using the "Basic Auth" protocol
func ProtectedRoute(ctn *dic.Container) func(c *fiber.Ctx) error {
	userPersistence := ctn.GetMysqlUserPersistence()
	logger := ctn.GetLogger()

	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			matched, _, err := userPersistence.BasicAuth(user, pass)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("error_protected_route")
				return false
			}
			if !matched {
				return false
			}
			return true
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(http.StatusUnauthorized).JSON(helpers.WrapStrInErrMap("Unauthorized"))
		},
		ContextUsername: user,
		ContextPassword: pass,
	})
}

func ExtractAuthedUserMeta(c *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(c)
	if err != nil {
		return err
	}
	logger := ctn.GetLogger()
	userPersistence := ctn.GetMysqlUserPersistence()

	user := c.Locals(user).(string)
	pass := c.Locals(pass).(string)

	_, userMeta, err := userPersistence.BasicAuth(user, pass)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("error_extract_authed_user_meta")
		return c.Status(http.StatusUnauthorized).JSON(helpers.WrapErrInErrMap(err))
	}
	c.Locals("userMeta", userMeta)
	return c.Next()
}
