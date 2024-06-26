package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
	"github.com/joho/godotenv"
)

func parseCSV(allowedItemsCSVPath string) (_ []functioncheck.AllowedItem, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error parsing %s: %w", allowedItemsCSVPath, err)
		}
	}()

	file, err := os.Open(allowedItemsCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var allowedItems []functioncheck.AllowedItem
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		allowedItem := functioncheck.AllowedItem{
			Name: line[0],
			Type: line[1],
		}
		allowedItems = append(allowedItems, allowedItem)
	}
	return allowedItems, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	allowedItems, _ := parseCSV("allowedItems.csv")
	err = functioncheck.Execute(allowedItems, "ex00")
	if err != nil {
		log.Print(err)

	}
	if err = git.Create("arthur"); err != nil {
		log.Fatalf("error: %s", err)
	}
}
