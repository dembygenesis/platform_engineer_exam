package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// initAPI boots our REST API connections
func initAPI(config *config.Config) {
	app := fiber.New(fiber.Config{
		BodyLimit: 20971520,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))

	app.Static("/", "./public")

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

	if err := app.Listen(":" + strconv.Itoa(config.API.Port)); err != nil {
		fmt.Println("gg has error", err)
	}
}

func main() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("error trying to initialize the dependency builder: %v", err.Error())
	}

	ctn := builder.Build()
	cfg, _ := ctn.SafeGetConfig()
	fmt.Println(cfg)

	initAPI(cfg)
}
