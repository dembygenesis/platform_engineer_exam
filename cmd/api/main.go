package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v", err.Error())
	}
	fmt.Println("cfg", cfg)
}
