package config

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/datastructures"
)

func GetAllowedItems(allowedItemsCSVPath string) (_ []datastructures.AllowedItem, err error) {
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

	var allowedItems []datastructures.AllowedItem
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		allowedItem := datastructures.AllowedItem{
			Name: line[0],
			Type: line[1],
		}
		allowedItems = append(allowedItems, allowedItem)
	}
	return allowedItems, nil
}
