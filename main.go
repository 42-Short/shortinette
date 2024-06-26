package main

import (
	"fmt"

	"github.com/42-Short/shortinette/pkg/git"
	"github.com/42-Short/shortinette/pkg/tester"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}
	if err = git.Create("shortinette-test"); err != nil {
		fmt.Println(err)
	}
	if err = git.AddCollaborator("shortinette-test", "shortinette-test", "push"); err != nil {
		fmt.Println(err)
	}
	if err = tester.Run("testconfig/R00.yaml"); err != nil {
		fmt.Println(err)
	}
}
