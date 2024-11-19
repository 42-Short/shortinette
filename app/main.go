package main

import (
	"os"

	"github.com/42-Short/shortinette/logger"
	"github.com/42-Short/shortinette/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error.Fatalf("cant load .env file: %v", err)
	}
}

func main() {
	server := server.NewServer(os.Getenv("SERVER_ADDR"))
	err := server.Run()
	if err != nil {
		logger.Error.Fatalf("failed to run Server: %v", err)
	}
	select {}
}
