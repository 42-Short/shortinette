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

	group.POST("/modules", InsertItemHandler(moduleDAO, api.timeout))
	group.POST("/participants", InsertItemHandler(participantDAO, api.timeout))

	group.PUT("/modules", UpdateItemHandler(moduleDAO, api.timeout))
	group.PUT("/participants", UpdateItemHandler(participantDAO, api.timeout))

	group.GET("/modules", GetAllItemsHandler(moduleDAO, api.timeout))
	group.GET("/participants", GetAllItemsHandler(participantDAO, api.timeout))

	group.GET("/modules/:id/:intra_login", GetItemHandler(moduleDAO, api.timeout))
	group.GET("/participants/:intra_login", GetItemHandler(participantDAO, api.timeout))

	group.DELETE("/modules/:id/:intra_login", DeleteItemHandler(moduleDAO, api.timeout))
	group.DELETE("/participants/:intra_login", DeleteItemHandler(participantDAO, api.timeout))
}
