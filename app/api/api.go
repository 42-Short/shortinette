package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

type API struct {
	*http.Server

	Engine *gin.Engine
	DB     *db.DB

	accessToken string
}

// Initializes and returns a new API instance
func NewAPI(addr string, db *db.DB, accessToken string, mode string) *API {
	engine := gin.Default()
	gin.SetMode(mode)

	return &API{
		Server: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
		Engine:      engine,
		DB:          db,
		accessToken: accessToken,
	}
}

// Starts the server in a go routine and listens for incoming requests.
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
