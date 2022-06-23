package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"net/http"
	"time"
)

func Throttle() func(c *fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 5 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(http.StatusForbidden).SendString("Please limit your requests to 5 per 5 seconds")
		},
	})
}
