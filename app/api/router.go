package api

import (
	"github.com/42-Short/shortinette/data"
)

func (api *API) setupRoutes() {
	group := api.Engine.Group("/shortinette/v1")
	group.Use(tokenAuthMiddleware())

	moduleDAO := data.NewDAO[data.Module](api.DB)
	participantDAO := data.NewDAO[data.Participant](api.DB)

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
