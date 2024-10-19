package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func Router() (router *gin.Engine) {
	router = gin.Default()

	router.GET("/admin", helloWorld)
	return router
}
