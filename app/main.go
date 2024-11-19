package main

import (
	"fmt"

	"github.com/42-Short/shortinette/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error.Fatalf("cant load .env file: %v", err)
	}
}

func main() {
	r := NewRouter()
	if err := r.Run("0.0.0.0:5000"); err != nil {
		fmt.Printf("error running gin server: %v", err)
	}
}
