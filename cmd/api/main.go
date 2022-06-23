package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/api"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func mapRoutes(app *fiber.App) {
	api := app.Group("/api", cors.New())
	apiToken := api.Group("/token")

	apiToken.Get("/")
}

func addContainerInstance(container *dic.Container) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals("dependencies", container)
		return c.Next()
	}
}

// initAPI boots our REST API connections
func initAPI(container *dic.Container, port string) {
	app := fiber.New(fiber.Config{
		BodyLimit: 20971520,
	})

	app.Use(addContainerInstance(container))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))

	app.Static("/", "./public")
	api.GetRouter(app)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		fmt.Println("Gracefully shutting down")
		err := app.Shutdown()
		if err != nil {
			fmt.Println("Shutting down error", err)
		}
	}()

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("error listening to port: %v, with msg: %v", port, err.Error())
	}
}

func main() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("error trying to initialize the dependency builder: %v", err.Error())
	}

	ctn := builder.Build()
	cfg, err := ctn.SafeGetConfig()
	if err != nil {
		log.Fatalf("error trying to fetch the config dependency: %v", err.Error())
	}

	initAPI(ctn, strconv.Itoa(cfg.API.Port))
}
