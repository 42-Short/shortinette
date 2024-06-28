package config

import (
	"fmt"
	"io"
	"os"

	"github.com/42-Short/shortinette/internal/datastructures"
	"gopkg.in/yaml.v2"
)

func GetAllowedItems(configFilePath string) (map[string][]datastructures.AllowedItem, error) {
	config, err := GetConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file %s: %w", configFilePath, err)
	}

	allowedItemsMap := make(map[string][]datastructures.AllowedItem)
	for key, exercise := range config.Exercises {
		var allowedItems []datastructures.AllowedItem
		for _, macro := range exercise.AllowedItems.Macros {
			allowedItems = append(allowedItems, datastructures.AllowedItem{
				Name: macro,
				Type: "macro",
			})
		}
		for _, function := range exercise.AllowedItems.Functions {
			allowedItems = append(allowedItems, datastructures.AllowedItem{
				Name: function,
				Type: "function",
			})
		}
		allowedItemsMap[key] = allowedItems
	}

	return allowedItemsMap, nil
}

func GetConfig(configFilePath string) (*datastructures.Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", configFilePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read from %s: %w", configFilePath, err)
	}

	var config datastructures.Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal yaml %s: %w", configFilePath, err)
	}

	return &config, nil
}
