package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/42-Short/shortinette/api"
	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error.Fatalf("can't load .env file: %v", err)
	}
}

func shutdown(api *api.API, sigCh chan os.Signal) {
	sig := <-sigCh
	err := api.Shutdown()
	if err != nil {
		logger.Error.Fatalf("failed to shutdown server: %v", err)
	}
	logger.Error.Fatalf("caught signal: %v", sig)
}

func run() {
	db, err := db.NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		logger.Error.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	err = db.Initialize("./db/schema.sql")
	if err != nil {
		logger.Error.Fatalf("failed to initialize db: %v", err)
	}

	_, err = data.SeedDB(db)
	if err != nil {
		logger.Error.Fatalf("failed to seed DB: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	api := api.NewAPI(db, gin.TestMode, time.Minute)
	go shutdown(api, sigCh)
	err = api.Run()
	if err != nil {
		logger.Error.Fatalf("failed to run api: %v", err)
	}
}

func main() {
	run()
	fmt.Printf("Hello World\n")
}
