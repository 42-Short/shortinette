package main

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/config"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}
	env, err := git.GetEnviorment();
	if err != nil {
		fmt.Println(err)
	}

	allowedItems, err := config.GetAllowedItems("testconfig/R00.yaml")
	if err != nil {
		fmt.Println(err)
	}
	if err = functioncheck.Execute(allowedItems, "ex00", env); err != nil {
		fmt.Println(err)
	}
	if err = git.Create("shortinette-test", env); err != nil {
		fmt.Println(err)
	}
	if err = git.AddCollaborator("shortinette-test", "shortinette-test", "push", env); err != nil {
		fmt.Println(err)
	}
}
