package middlewares

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	userKey = "user"
	passKey = "pass"
)

// ProtectedRoute guards a route using the "Basic Auth" protocol
func ProtectedRoute(ctn *dic.Container) func(ctx *fiber.Ctx) error {
	userPersistence := ctn.GetMysqlUserPersistence()

	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			matched, _, err := userPersistence.BasicAuth(user, pass)
			if err != nil {
				return false
			}
			if !matched {
				return false
			}
			return true
		},
		Unauthorized: func(ctx *fiber.Ctx) error {
			logger := common.GetLogger(ctx.Context())
			logger.WithFields(logrus.Fields{
				"msg": "Unauthorized",
			}).Error("error_protected_route")
			return ctx.Status(http.StatusUnauthorized).JSON(helpers.WrapStrInErrMap("Unauthorized"))
		},
		ContextUsername: userKey,
		ContextPassword: passKey,
	})
}

func ExtractAuthedUserMeta(ctx *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(ctx)
	if err != nil {
		return err
	}
	logger := common.GetLogger(ctx.Context())
	userPersistence := ctn.GetMysqlUserPersistence()

	user := ctx.Locals(userKey).(string)
	pass := ctx.Locals(passKey).(string)

	_, userMeta, err := userPersistence.BasicAuth(user, pass)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("error_extract_authed_user_meta")
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.WrapErrInErrMap(err))
	}
	ctx.Locals("userMeta", userMeta)
	return ctx.Next()
}
