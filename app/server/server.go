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
	Handler *gin.Engine
}

func NewServer(addr string) *Server {
	handler := gin.Default()
	gin.SetMode(gin.DebugMode)

	server := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		Handler: handler,
	}
	server.setupRoutes()
	return server
}

func (s *Server) Run() chan error {
	errorCh := make(chan error, 1)

	go func() {
		err := s.ListenAndServe()
		errorCh <- err
	}()
	logger.Info.Printf("server is listening on %s", s.Addr)
	return errorCh
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("faild to shut down Server with addr: %s: %v", s.Addr, err)
	}
	logger.Info.Printf("Server with addr: `%s` has  gracefully shut down.\n", s.Addr)
	return nil
}
