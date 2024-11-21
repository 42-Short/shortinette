package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

type API struct {
	*http.Server

	Engine      *gin.Engine
	DB          *db.DB
	timeout     time.Duration
	accessToken string
}

// Initializes and returns a new API instance
// timeout specifies the max duration per request
func NewAPI(db *db.DB, mode string, timeout time.Duration) *API {
	accessToken := os.Getenv("API_TOKEN")
	if accessToken == "" && mode != gin.TestMode {
		panic("API_TOKEN not found in .env")
	} else if mode == gin.TestMode {
		accessToken = "test"
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "localhost:8080"
	}

	engine := gin.Default()
	gin.SetMode(mode)
	api := &API{
		Server: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
		Engine:      engine,
		DB:          db,
		timeout:     timeout,
		accessToken: accessToken,
	}
	api.setupRoutes()
	return api
}

// Starts the server in a go routine and listens for incoming requests.
// returns a channel for error monitoring
func (api *API) Run() error {
	logger.Info.Printf("server is listening on %s", api.Addr)

	err := api.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to listen to %s", api.Addr)
	}
	return nil
}

// Gracefully shuts down the API server
func (api *API) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := api.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("faild to shut down Server with addr: %s: %v", api.Addr, err)
	}
	logger.Info.Printf("Server with addr: `%s` has  gracefully shut down.\n", api.Addr)
	return nil
}
