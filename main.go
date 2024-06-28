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
