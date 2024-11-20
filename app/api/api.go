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

func NewAPI(db *db.DB, mode string, timeout time.Duration) *API {
	accessToken := os.Getenv("API_TOKEN")
	if accessToken == "" {
		logger.Warning.Printf("API_TOKEN not found. Creating API without access Token")
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

func (api *API) Run() chan error {
	errorCh := make(chan error, 1)

	go func() {
		err := api.ListenAndServe()
		errorCh <- err
	}()
	logger.Info.Printf("server is listening on %s", api.Addr)

	return errorCh
}

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
