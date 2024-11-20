package api

import (
	"net/http"
	"os"

	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

func tokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("API_TOKEN")

	if requiredToken == "" {
		logger.Error.Fatal("API_TOKEN in .env is empty\n")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing Authorization header format"})
			c.Abort()
			return
		}

		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		if token != requiredToken {
			logger.Warning.Printf("Unauthorized access attempt with token: " + token + "\n")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "token invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
