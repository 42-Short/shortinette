package main

import (
	"log"

	// "github.com/42-Short/shortinette/cmd"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = git.Execute("https://github.com/42-student-council/website.git", "website")
	if err != nil {
		log.Fatal(err)
	}

	err = functioncheck.Execute("allowedItems.csv")
	if err != nil {
		log.Fatal(err)
	}
}
