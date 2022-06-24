package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/api"
	"github.com/dembygenesis/platform_engineer_exam/api/helpers"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
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
func initAPI(ctn *dic.Container, port string) {
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

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("error listening to port: %v, with msg: %v", port, err.Error())
	}
}

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

	mysql, err := ctn.SafeGetMysqlConnection()
	if err != nil {
		log.Fatalf("error trying to fetch the mysql connection from the container: %v", err.Error())
	}
	err = mysql.Ping()
	if err != nil {
		log.Fatalf("error trying to ping from the mysql persistence: %v", err.Error())
	}

	initAPI(ctn, strconv.Itoa(cfg.API.Port))
}
