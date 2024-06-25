package main

import (
	"log"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = functioncheck.Execute("allowedItems.csv")
	if err != nil {
		log.Fatal(err)
	}
}
