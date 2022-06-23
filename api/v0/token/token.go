package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetToken(c *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	userMeta, ok := c.Locals("userMeta").(*models_schema.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapStrInErrMap("userMeta conversion fails"))
	}

	generatedToken, err := biz.Generate(userMeta)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	return c.Status(http.StatusCreated).JSON(generatedToken)
}
