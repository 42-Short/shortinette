package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	r := server.NewRouter()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := r.Run("0.0.0.0:5000"); err != nil {
			fmt.Printf("error running gin server: %v", err)
		}
	}()
	<-done
}
