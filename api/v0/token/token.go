package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetToken(ctx *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	generatedToken, err := biz.Generate()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return ctx.JSON(generatedToken)
}
