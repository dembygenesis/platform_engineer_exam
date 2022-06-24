package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ValidateToken(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapStrInErrMap("token is missing"))
	}

	ctn, err := helpers.GetContainer(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	cfg, err := ctn.SafeGetConfig()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	err = biz.Validate(c.Context(), token, cfg.App.TokenLapseDuration, cfg.App.TokenLapseSettings)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	return c.Status(http.StatusOK).JSON(true)
}

func GetAll(c *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	tokens, err := biz.GetAll(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	return c.Status(http.StatusCreated).JSON(tokens)
}

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

	generatedToken, err := biz.Generate(c.Context(), userMeta)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	return c.Status(http.StatusCreated).JSON(generatedToken)
}
