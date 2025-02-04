package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

func githubAuthMiddleware(accessToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}
		c.Set(gin.BodyBytesKey, body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		signature := c.GetHeader("X-Hub-Signature-256")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication header for GitHub webhook"})
			c.Abort()
			return
		}
		fmt.Println(signature)
		mac := hmac.New(sha256.New, []byte(accessToken))
		mac.Write(body)
		expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication tokens do not match"})
			c.Abort()
			return
		}

		c.Next()
	}
}

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
