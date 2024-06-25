package functioncheck

import (
	"encoding/csv"
	"fmt"
	"os"
)

type AllowedItem struct {
	Name string
	Type string
}

func parseCSV(allowedItemsCSVPath string) ([]AllowedItem, error) {
	file, err := os.Open(allowedItemsCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var allowedItems []AllowedItem
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		allowedItem := AllowedItem{
			Name: line[0],
			Type: line[1],
		}
		allowedItems = append(allowedItems, allowedItem)
	}
	return allowedItems, nil
}

func Execute(allowedItemsCSVPath string) error {
	allowedItems, err := parseCSV(allowedItemsCSVPath)
	if err != nil {
		return fmt.Errorf("error parsing %s: %s", allowedItemsCSVPath, err)
	}
	fmt.Println(allowedItems)
	return nil
}
