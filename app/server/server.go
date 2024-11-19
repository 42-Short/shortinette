package server

import (
	"net/http"
	"os"

	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("API_TOKEN")

	if requiredToken == "" {
		logger.Error.Fatal("api access token is empty")
	}

	return func(c *gin.Context) {
		token := c.Request.FormValue("api_token")
		if token != requiredToken {
			c.JSON(http.StatusBadRequest, gin.H{"message": "token invalid"})
			return
		}

		c.Next()
	}

}

func NewRouter() *gin.Engine {
	r := gin.Default() //TODO: check if options are required

	group := r.Group("v1/")
	group.Use(TokenAuthMiddleware())

	return r
}
