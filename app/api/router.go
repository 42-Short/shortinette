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

	group.POST("/modules", insertItemHandler(moduleDAO, api.timeout))
	group.POST("/participants", insertItemHandler(participantDAO, api.timeout))

	group.PUT("/modules", updateItemHandler(moduleDAO, api.timeout))
	group.PUT("/participants", updateItemHandler(participantDAO, api.timeout))

	group.GET("/modules", getAllItemsHandler(moduleDAO, api.timeout))
	group.GET("/participants", getAllItemsHandler(participantDAO, api.timeout))

	group.GET("/modules/:id/:intra_login", getItemHandler(moduleDAO, api.timeout))
	group.GET("/participants/:intra_login", getItemHandler(participantDAO, api.timeout))

	group.DELETE("/modules/:id/:intra_login", deleteItemHandler(moduleDAO, api.timeout))
	group.DELETE("/participants/:intra_login", deleteItemHandler(participantDAO, api.timeout))

	group.POST("/webhook", githubWebhookHandler())
}
