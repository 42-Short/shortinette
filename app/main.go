package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/42-Short/shortinette/api"
	"github.com/42-Short/shortinette/config"
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
	ex0, _ := config.NewExercise(
		10,
		[]string{"hello.rs"},
		"ex00",
	)
	ex1, _ := config.NewExercise(
		10,
		[]string{"min.rs"},
		"ex01",
	)
	ex2, _ := config.NewExercise(
		10,
		[]string{"yes.rs", "collatz.rs", "print_bytes.rs"},
		"ex02",
	)
	ex3, _ := config.NewExercise(
		10,
		[]string{"fizzbuzz.rs"},
		"ex03",
	)
	ex4, _ := config.NewExercise(
		10,
		[]string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"},
		"ex04",
	)
	ex5, _ := config.NewExercise(
		15,
		[]string{"src/main.rs", "src/lib.rs", "Cargo.toml"},
		"ex05",
	)
	ex6, _ := config.NewExercise(
		15,
		[]string{"src/main.rs", "Cargo.toml"},
		"ex06",
	)
	ex7, _ := config.NewExercise(
		20,
		[]string{"src/lib.rs", "src/main.rs", "Cargo.toml"},
		"ex02",
	)

	module0, _ := config.NewModule(
		[]config.Exercise{*ex0, *ex1, *ex2, *ex3, *ex4, *ex5, *ex6, *ex7},
		15,
	)

	ex0, _ = config.NewExercise(
		10,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex00",
	)
	ex1, _ = config.NewExercise(
		10,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex01",
	)
	ex2, _ = config.NewExercise(
		10,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex02",
	)
	ex3, _ = config.NewExercise(
		10,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex03",
	)
	ex4, _ = config.NewExercise(
		10,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex04",
	)
	ex5, _ = config.NewExercise(
		15,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex05",
	)
	ex6, _ = config.NewExercise(
		15,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex06",
	)
	ex7, _ = config.NewExercise(
		20,
		[]string{"src/lib.rs", "Cargo.toml"},
		"ex02",
	)

	module1, _ := config.NewModule(
		[]config.Exercise{*ex0, *ex1, *ex2, *ex3, *ex4, *ex5, *ex6, *ex7},
		15,
	)

	return config.NewConfig(
		[]config.Module{*module0, *module1},
		24*time.Hour,
		time.Now(),
		"42short/rust",
		"./rust",
	)
}

func run() {
	db, err := db.NewDB(context.Background(), "./data/shortinette.db")
	if err != nil {
		logger.Error.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	err = db.Initialize("./db/schema.sql")
	if err != nil {
		logger.Error.Fatalf("failed to initialize db: %v", err)
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
