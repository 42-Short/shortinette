package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func router() (router *gin.Engine) {
	router = gin.Default()

	router.GET("/admin", helloWorld)

	return router
}

func main() {
	r := router()
	r.Run("0.0.0.0:5000")
}
