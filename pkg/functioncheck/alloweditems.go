package functioncheck

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/internal/datastructures"
)

func writeTemplateToFile(template, itemName string, file *os.File) error {
	content := fmt.Sprintf(template, itemName, itemName)
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error writing entry: %w", err)
	}
	return nil
}

func writeAllowedItemsLib(allowedItems []datastructures.AllowedItem, file *os.File) error {
	exerciseNumber := "00"
	content := fmt.Sprintf(allowedItemsLibHeader, exerciseNumber)
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	for _, item := range allowedItems {
		if item.Type == "macro" {
			if err := writeTemplateToFile(allowedMacroTemplate, item.Name, file); err != nil {
				return err
			}
		} else if item.Type == "function" {
			if err := writeTemplateToFile(allowedFunctionTemplate, item.Name, file); err != nil {
				return err
			}
		}
	}
	if _, err := file.WriteString("}"); err != nil {
		return fmt.Errorf("error writing closing bracket: %w", err)
	}
	return nil
}

func createFileWithDirs(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func writeCargoToml(filePath, content string) error {
	file, err := createFileWithDirs(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}
