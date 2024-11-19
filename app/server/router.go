package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() {
	group := s.Handler.Group("/shortinette/v1")
	group.Use(tokenAuthMiddleware())

	group.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
}
