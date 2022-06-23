package helpers

import (
	"errors"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrNoContainerFound = errors.New("no container found in the fiber context provided")
)

func GetContainer(ctx *fiber.Ctx) (*dic.Container, error) {
	ctn, ok := ctx.Locals(Dependencies).(*dic.Container)
	if !ok {
		return nil, ErrNoContainerFound
	}
	return ctn, nil
}

// AddContainerInstance injects our dependencies to our fiber context
func AddContainerInstance(container *dic.Container) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(Dependencies, container)
		return c.Next()
	}
}

func WrapInErrMap(i interface{}) interface{} {
	w := make(map[string][]interface{})
	w["errors"] = []interface{}{i}
	return w
}
