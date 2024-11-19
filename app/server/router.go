package server

import (
	"net/http"

	"github.com/42-Short/shortinette/data"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() {
	group := s.Handler.Group("/shortinette/v1")
	group.Use(tokenAuthMiddleware())

	moduleDAO := data.NewDAO[data.Module](nil)           //TODO
	participantDAO := data.NewDAO[data.Participant](nil) //TODO

	group.Any("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	group.POST("/modules", InsertItemHandler(moduleDAO))
	group.POST("/participants", InsertItemHandler(participantDAO))

	group.PUT("/modules/:intra_login/:id", UpdateItemHandler(moduleDAO))
	group.PUT("/participants/:intra_login", UpdateItemHandler(participantDAO))

	group.GET("/modules", GetAllItemsHandler(moduleDAO))
	group.GET("/participants", GetAllItemsHandler(participantDAO))

	group.GET("/modules/:intra_login/:id", GetItemHandler(moduleDAO))
	group.GET("/participants/:intra_login", GetItemHandler(participantDAO))

	group.DELETE("/modules/:intra_login/:id", DeleteItemHandler(moduleDAO))
	group.DELETE("/participants/:intra_login", DeleteItemHandler(participantDAO))

}
