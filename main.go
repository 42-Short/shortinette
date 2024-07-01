package main

import (
	"log"

	"github.com/42-Short/shortinette/internal/endpoints"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/utils"
)

func main() {
	logger.InitializeStandardLoggers()
	if err := utils.RequireEnv(); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	if err := endpoints.CreateNewTeam("shortinette-test", "R00"); err != nil {
		log.Fatalf("could not create team: %s", err)
	}
	if result, err := endpoints.TestSubmission("shortinette-test-R00", "testconfig/R00.yaml"); err != nil {
		logger.Error.Println(err)
	} else {
		logger.Info.Printf("tests run successfully, results: %s", result)
	}
}
