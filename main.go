package main

import (
	"log"

	"github.com/42-Short/shortinette/internal/config"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	allowedItems, _ := config.GetAllowedItems("allowedItems.csv")
	err = functioncheck.Execute(allowedItems, "ex00")
	if err != nil {
		log.Print(err)

	}
	if err = git.Create("shortinette-test"); err != nil {
		log.Printf("error: %s", err)
	}
	if err = git.AddCollaborator("shortinette-test", "shortinette-test", "push"); err != nil {
		log.Printf("error: %s", err)
	}
}
