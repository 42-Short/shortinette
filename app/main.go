package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/42-Short/shortinette/api"
	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

func shutdown(api *api.API, sigCh chan os.Signal) {
	sig := <-sigCh
	err := api.Shutdown()
	if err != nil {
		logger.Error.Fatalf("failed to shutdown server: %v", err)
	}
	logger.Error.Fatalf("caught signal: %v", sig)
}

func getMockConfig() *config.Config {
	ex1, _ := config.NewExercise(
		10,
		[]string{"*.c", "*.h"},
		"ex00",
	)
	ex2, _ := config.NewExercise(
		20,
		[]string{"*.c", "*.h"},
		"ex01",
	)

	module, _ := config.NewModule(
		[]config.Exercise{*ex1, *ex2},
		15,
	)

	return config.NewConfig(
		[]config.Module{*module, *module},
		24*time.Hour,
		time.Now(),
		"app/testenv/test.sh",
	)
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

	config := getMockConfig()

	if err := config.FetchEnvVariables(); err != nil {
		logger.Error.Fatalf("could not fetch environment variables: %v", err)
	}

	api := api.NewAPI(config, db, gin.DebugMode)
	api.SetupRouter()
	go shutdown(api, sigCh)
	err = api.Run()
	if err != nil {
		logger.Error.Fatalf("failed to run api: %v", err)
	}
}

func main() {
	run()
}
