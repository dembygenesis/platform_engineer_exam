package token

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetToken(ctx *fiber.Ctx) error {
	ctn := helpers.GetContainer(ctx)
	cfg, err := ctn.SafeGetConfig()
	if err != nil {
		panic("error getting config: " + err.Error())
	}
	fmt.Println("cfg", cfg)
	return ctx.JSON("Test")
}
