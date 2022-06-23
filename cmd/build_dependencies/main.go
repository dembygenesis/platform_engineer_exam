package main

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/provider"
	"github.com/sarulabs/dingo/v4"
	"os"
)

func main() {
	// Compile
	if len(os.Args) != 2 {
		fmt.Println("usage: go run main.go ~/apps/platform_engineer_exam/dependency_injection")
		os.Exit(1)
	}

	err := dingo.GenerateContainer((*provider.Provider)(nil), os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
