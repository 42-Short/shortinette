package main

import (
	"fmt"
	"log"

	// "github.com/42-Short/shortinette/cmd"
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
		fmt.Println(err)
	}
}
