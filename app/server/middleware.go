package server

import (
	"net/http"
	"os"

	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

func tokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("API_TOKEN")

	if requiredToken == "" {
		logger.Error.Fatal("API_TOKEN in .env is empty")
	}

	return func(c *gin.Context) {
		token := c.Request.FormValue("api_token")
		if token != requiredToken {
			c.JSON(http.StatusBadRequest, gin.H{"message": "token invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
