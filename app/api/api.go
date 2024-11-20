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
}

func NewAPI(addr string, db *db.DB) *API {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)
	api := &API{
		Server: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
		Engine: engine,
		DB:     db,
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
