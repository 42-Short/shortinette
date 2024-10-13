//go:build ignore

package main

import (
	"os"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/webserver"
	"github.com/42-Short/shortinette/rust/tests/R00"
	"github.com/42-Short/shortinette/rust/tests/R01"
	"github.com/42-Short/shortinette/rust/tests/R02"
	"github.com/42-Short/shortinette/rust/tests/R03"
	"github.com/42-Short/shortinette/rust/tests/R04"
	"github.com/42-Short/shortinette/rust/tests/R05"
	"github.com/42-Short/shortinette/rust/tests/R06"
	"github.com/gin-gonic/gin"
)

var modules = map[string]Module.Module{
	"00": *R00.R00(),
	"01": *R01.R01(),
	"02": *R02.R02(),
	"03": *R03.R03(),
	"04": *R04.R04(),
	"05": *R05.R05(),
	"06": *R06.R06(),
}

func router() *gin.Engine {
	router := gin.Default()

	webhook := webserver.NewWebhook(modules)
	router.POST("/webhook", webhook.HandleWebhook)

	return router
}

func main() {
	logger.InitializeStandardLoggers("")

	router := router()
	if err := router.Run(os.Getenv("WEBHOOK_PORT")); err != nil {
		logger.Error.Printf("could not start gin router: %v", err)
	}
}
