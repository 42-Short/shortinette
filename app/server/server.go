package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
}

func NewServer(addr string) *Server {
	handler := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	group := handler.Group("v1/")
	group.Use(tokenAuthMiddleware())

	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (server *Server) Run() error {
	// go func (){
	// 	err := server.ListenAndServe()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	logger.Info.Printf("Server is listening at %s...", server.Addr)
	return nil
}

func (server *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("faild to shut down Server with addr: %s gracefully: %v", server.Addr, err)
	}
	logger.Info.Printf("Server has gracefully shut down for addr: %s\n", server.Addr)
	return nil
}
