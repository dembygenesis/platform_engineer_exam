package api

import (
	"github.com/dembygenesis/platform_engineer_exam/api/v0/middlewares"
	"github.com/dembygenesis/platform_engineer_exam/api/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
)

func GetRouter(app *fiber.App, ctn *dic.Container) {
	api := app.Group("/api")
	v0 := api.Group("/v0")

	v0token := v0.Group("/token")
	v0token.Get("/", middlewares.ProtectedRoute(ctn), middlewares.ExtractAuthedUserMeta, token.GetAll)
	v0token.Post("/", middlewares.ProtectedRoute(ctn), middlewares.ExtractAuthedUserMeta, token.GetToken)
	v0token.Get("/:token/validate", middlewares.Throttle(), token.ValidateToken)
}
