package api

import (
	"net/http"

	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

func tokenAuthMiddleware(accessToken string) gin.HandlerFunc {
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

		if token != accessToken {
			logger.Warning.Printf("unauthorized access attempt with token: %s \n", token)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "token invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
