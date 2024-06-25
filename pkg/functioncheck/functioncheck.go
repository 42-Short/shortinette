package functioncheck

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
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

func initCompilingEnvironment(allowedItems []AllowedItem, ) error {
	libFilePath := "functioncheck/allowedfunctions/src/lib.rs"
	dir := filepath.Dir(libFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(libFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	exerciseNumber := "00"
	content := fmt.Sprintf("pub mod ex%s { ", exerciseNumber)
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	for _, item := range allowedItems {
		if item.Type == "macro" {
			content := fmt.Sprintf(`
		#[cfg(not(feature = "allowed_%[1]s"))]
		#[macro_export]
		macro_rules! %[1]s {
			($($arg:tt)*) => {{}}
	}
`, item.Name)
			if _, err := file.WriteString(content); err != nil {
				return err
			}
		}
	}
	file.WriteString("}")
	return nil
}

func Execute(allowedItemsCSVPath string) error {
	allowedItems, err := parseCSV(allowedItemsCSVPath)
	if err != nil {
		return fmt.Errorf("error parsing %s: %s", allowedItemsCSVPath, err)
	}
	err = initCompilingEnvironment(allowedItems)
	if err != nil {
		return fmt.Errorf("error initializing compiling environment: %s", err)
	}
	return nil
}
