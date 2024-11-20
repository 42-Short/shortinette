package api

import (
	"github.com/42-Short/shortinette/data"
)

func (api *API) setupRoutes() {
	group := api.Engine.Group("/shortinette/v1")
	if api.accessToken != "" {
		group.Use(tokenAuthMiddleware(api.accessToken))
	}

	moduleDAO := data.NewDAO[data.Module](api.DB)
	participantDAO := data.NewDAO[data.Participant](api.DB)

	group.POST("/modules", InsertItemHandler(moduleDAO))
	group.POST("/participants", InsertItemHandler(participantDAO))

	group.PUT("/modules", UpdateItemHandler(moduleDAO)) //TODO: add id and intra login to update
	group.PUT("/participants", UpdateItemHandler(participantDAO))

	group.GET("/modules", GetAllItemsHandler(moduleDAO))
	group.GET("/participants", GetAllItemsHandler(participantDAO))

	group.GET("/modules/:id/:intra_login", GetItemHandler(moduleDAO))
	group.GET("/participants/:intra_login", GetItemHandler(participantDAO))

	group.DELETE("/modules/:id/:intra_login", DeleteItemHandler(moduleDAO))
	group.DELETE("/participants/:intra_login", DeleteItemHandler(participantDAO))
}
