package main

import (
	"context"
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
		logger.Error.Fatalf("can't load .env file: %v", err)
	}
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandler(cancel)

	server := server.NewServer(os.Getenv("SERVER_ADDR"))
	errCh := server.Run()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				logger.Error.Printf("failed to run server: %v", err)
			}
			os.Exit(1)
		case <-ctx.Done():
			if err := server.Shutdown(); err != nil {
				logger.Error.Printf("failed to shutdown server: %v", err)
			}
			os.Exit(0)
		default:
			// do stuff
		}
	}
}
