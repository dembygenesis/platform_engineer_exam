package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"log"
)

func main() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("error trying to initialize the dependency builder: %v", err.Error())
	}

	ctn := builder.Build()
	cfg, _ := ctn.SafeGetConfig()
	fmt.Println(cfg)
}
