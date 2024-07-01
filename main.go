package main

import (
	"fmt"
	"log"

	"github.com/42-Short/shortinette/internal/endpoints"
	"github.com/42-Short/shortinette/internal/utils"
)

func main() {
	if err := utils.RequireEnv(); err != nil {
		log.Fatalf(err.Error())
	}
	if err := endpoints.CreateNewTeam("shortinette-test", "R00"); err != nil {
		log.Fatalf("could not create team: %s", err)
	}
	if result, err := endpoints.TestSubmission("shortinette-test-R00", "testconfig/R00.yaml"); err != nil {
		log.Fatalf("could not run tests: %s", err)
	} else {
		fmt.Printf("tests run successfully: %s", result)
	}
}
