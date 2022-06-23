package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetToken(c *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	generatedToken, err := biz.Generate()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.JSON(generatedToken)
}
