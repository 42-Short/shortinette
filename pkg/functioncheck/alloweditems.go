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

func writeAllowedItemsLib(allowedItems datastructures.AllowedItems, file *os.File, exercise string) error {
	content := fmt.Sprintf(allowedItemsLibHeader, exercise)
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	for _, macro := range allowedItems.Macros {
		if err := writeTemplateToFile(allowedMacroTemplate, macro, file); err != nil {
			return err
		}
	}
	for _, function := range allowedItems.Functions {
		if err := writeTemplateToFile(allowedFunctionTemplate, function, file); err != nil {
			return err
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
