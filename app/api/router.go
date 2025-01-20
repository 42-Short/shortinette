package api

import (
	"github.com/42-Short/shortinette/data"
)

func (api *API) SetupRouter() {
	group := api.Engine.Group("/shortinette/v1")
	group.Use(tokenAuthMiddleware(api.config.ApiToken))

	moduleDAO := data.NewDAO[data.Module](api.DB)
	participantDAO := data.NewDAO[data.Participant](api.DB)

	group.POST("/webhook/grademe", githubWebhookHandler(moduleDAO, participantDAO, *api.config))
	group.Any("/modules/:id/:intra_login/grademe", gradingHandler(moduleDAO, participantDAO, *api.config))

	group.POST("/modules", insertItemHandler(moduleDAO))
	group.POST("/participants", insertItemHandler(participantDAO))

	group.PUT("/modules", updateItemHandler(moduleDAO))
	group.PUT("/participants", updateItemHandler(participantDAO))

	group.GET("/modules", getAllItemsHandler(moduleDAO))
	group.GET("/participants", getAllItemsHandler(participantDAO))

	group.GET("/modules/:id/:intra_login", getItemHandler(moduleDAO))
	group.GET("/participants/:intra_login", getItemHandler(participantDAO))

	group.DELETE("/modules/:id/:intra_login", deleteItemHandler(moduleDAO))
	group.DELETE("/participants/:intra_login", deleteItemHandler(participantDAO))

}
