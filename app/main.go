package main

import (
	"fmt"

	"github.com/42-Short/shortinette/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error.Fatalf("can't load .env file: %v", err)
	}
}

func main() {
	fmt.Printf("Hello World\n")
}
