package api

import (
	"github.com/dembygenesis/platform_engineer_exam/api/v0/middlewares"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
)

func GetRouter(app *fiber.App, ctn *dic.Container) {
	api := app.Group("/api")
	v0 := api.Group("/v0")

	// apiToken := token.NewAPIToken()
	apiToken := ctn.GetApiToken()
	authMiddlewares := ctn.GetApiMiddlewares()

	v0token := v0.Group("/token")
	v0token.Get("/", authMiddlewares.ProtectedRoute(), authMiddlewares.AttachUserMeta, apiToken.GetAll)
	v0token.Post("/", authMiddlewares.ProtectedRoute(), authMiddlewares.AttachUserMeta, apiToken.GetToken)
	v0token.Get("/:token/validate", middlewares.Throttle(), apiToken.ValidateToken)
	v0token.Delete("/:token/revoke", authMiddlewares.ProtectedRoute(), authMiddlewares.AttachUserMeta, apiToken.Revoke)
}
