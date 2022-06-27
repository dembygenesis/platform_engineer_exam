package token

import (
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// ValidateToken
// @Id ValidateToken
// @Summary Validate
// @Description Validates a string token passed
// @Tags Token
// @Param token path string true "token"
// @Accept application/json
// @Produce application/json
// @Success 200 {boolean} boolean
// @Failure 400 {object} models.AuthFailBadRequest
// @Failure 500 {object} models.AuthFailInternalServerError
// @Router /v0/token/{token}/validate [get]
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

// GetAll
// @Id GetAll
// @Summary Fetch all
// @Description Fetches all tokens added by the admin user
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

// Revoke
// @Id Revoke
// @Summary Revoke
// @Description Revokes a token's access
// @Tags Token
// @Accept application/json
// @Produce application/json
// @Param token path string true "token"
// @Success 200 {boolean} boolean
// @Failure 400 {object} models.AuthFailBadRequest
// @Failure 500 {object} models.AuthFailInternalServerError
// @Security BasicAuth
// @Router /v0/token/{token}/revoke [delete]
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
// @Summary Create
// @Description Creates a new invite token
// @Tags Token
// @Accept application/json
// @Produce application/json
// @Success 201 {string} string
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
