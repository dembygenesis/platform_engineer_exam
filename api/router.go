package api

import (
	"github.com/dembygenesis/platform_engineer_exam/api/v0/middlewares/auth"
	"github.com/dembygenesis/platform_engineer_exam/api/v0/token"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
)

func GetRouter(app *fiber.App, ctn *dic.Container) {
	api := app.Group("/api")
	v0 := api.Group("/v0")

	v0token := v0.Group("/token")
	v0token.Get("/", auth.ProtectedRoute(ctn), token.GetToken)
}
