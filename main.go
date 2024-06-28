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
	env, err := git.GetEnvironment();
	if err != nil {
		fmt.Println(err)
	}
	if err = tester.Run("testconfig/R00.yaml", "shortinette-test", "studentcode", env); err != nil {
		fmt.Println(err)
	}
}
