package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/api"
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// initAPI boots our REST API connections
func initAPI(ctn *dic.Container, cfg *config.Config) {
	app := fiber.New(fiber.Config{
		BodyLimit: 20971520,
	})

	app.Use(requestid.New())
	app.Use(helpers.AddContainerInstance(ctn))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))

	// app.Static("/docs", "./docs/index.html")
	app.Static("/docs", "./docs")
	app.Static("/", "./public")
	api.GetRouter(app, ctn)

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

	port := strconv.Itoa(cfg.API.Port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("error listening to port: %v, with msg: %v", port, err.Error())
	}
}

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath http://localhost:8081/api
// @securityDefinitions.basic BasicAuth
// @in header
// @name HeaderAuth
func main() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("error trying to initialize the builder: %v", err.Error())
	}
	ctn := builder.Build()

	cfg, err := ctn.SafeGetConfig()
	if err != nil {
		log.Fatalf("error trying to fetch the config from the container: %v", err.Error())
	}

	initAPI(ctn, cfg)
}
