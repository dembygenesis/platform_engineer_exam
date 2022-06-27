package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ValidateToken(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	if token == "" {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapStrInErrMap("token is missing"))
	}

	ctn, err := helpers.GetContainer(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	err = biz.Validate(ctx.Context(), token)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	return ctx.Status(http.StatusOK).JSON(true)
}

// GetAll func for creates a new book.
// @Id GetAll
// @Description Fetches all tokens added by the admin user
// @Summary Fetches all tokens
// @Tags Token
// @Accept application/json
// @Produce application/json
// @Success 200 {object} []models.Token
// @Failure 400 {object} models.AuthFailBadRequest
// @Failure 500 {object} models.AuthFailInternalServerError
// @Security BasicAuth
// @Router /v0/token [get]
func GetAll(ctx *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	tokens, err := biz.GetAll(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	return ctx.Status(http.StatusCreated).JSON(tokens)
}

func Revoke(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	if token == "" {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapStrInErrMap("token is missing"))
	}

	ctn, err := helpers.GetContainer(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	err = biz.Revoke(ctx.Context(), token)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	return ctx.Status(http.StatusOK).SendString("Revoked token access!")
}

// GetToken Creates a new invite token
// @Id GetToken
// @Description Creates a new invite token
// @Summary Creates a new invite token
// @Tags Token
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string
// @Failure 400 {object} models.AuthFailBadRequest
// @Failure 500 {object} models.AuthFailInternalServerError
// @Security BasicAuth
// @Router /v0/token [post]
func GetToken(ctx *fiber.Ctx) error {
	ctn, err := helpers.GetContainer(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}
	biz, err := ctn.SafeGetBusinessToken()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	userMeta, ok := ctx.Locals("userMeta").(*models.User)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapStrInErrMap("userMeta conversion fails"))
	}

	generatedToken, err := biz.Generate(ctx.Context(), userMeta)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.WrapErrInErrMap(err))
	}

	return ctx.Status(http.StatusCreated).JSON(generatedToken)
}
