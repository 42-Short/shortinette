package main

import (
	"fmt"

	"github.com/42-Short/shortinette/pkg/git"
	"github.com/42-Short/shortinette/pkg/tester"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}
	err := git.CheckRequiredEnvironmentVariables();
	if err != nil {
		fmt.Println(err)
	}
	if err = tester.Run("testconfig/R00.yaml", "shortinette-test", "studentcode"); err != nil {
		fmt.Println(err)
	}
}
