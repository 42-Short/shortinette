package main

import (
	"net/http"

	"github.com/42-Short/shortinette/tester"
	"github.com/42-Short/shortinette/logger"
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
	done := make(chan bool, 1)

	tester.HandleSignals(done, true)
	r := router()
	if err := r.Run("0.0.0.0:5000"); err != nil {
		logger.Error.Printf("error running gin server: %v", err)
	}
	done <- true

	<-done
}
