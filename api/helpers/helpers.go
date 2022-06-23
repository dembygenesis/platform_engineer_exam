package helpers

import (
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
)

func GetContainer(ctx *fiber.Ctx) *dic.Container {
	return ctx.Locals("dependencies").(*dic.Container)
}
