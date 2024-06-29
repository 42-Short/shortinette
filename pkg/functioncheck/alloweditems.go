package functioncheck

import (
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"strings"

	"github.com/42-Short/shortinette/internal/datastructures"
	"github.com/42-Short/shortinette/internal/templates"
)

func searchForKeyword(keywords []datastructures.Keyword, word string) (keyword datastructures.Keyword, found bool) {
	for _, keyword := range keywords {
		if word == keyword.Keyword{
			return keyword, true
		}
	}
	return keyword, false
}

func checkKeywordAmount(keywordCounts map[string] int,  keywords []datastructures.Keyword) (err error) {
	foundKeywords := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
        if count, inMap := keywordCounts[keyword.Keyword]; inMap {
            if count > keyword.Amount {
				foundKeywords = append(foundKeywords, keyword.Keyword)
            }
        }
    }
	if len(foundKeywords) > 0 {
		return fmt.Errorf("keywords %s are used more often than allowed", strings.Join(foundKeywords, ", "))
	}
	return nil
}

func scanStudentFile(scanner *bufio.Scanner, allowedItems datastructures.AllowedItems) (err error) {
	keywordCounts := make(map[string] int)
	for scanner.Scan() {
        word := scanner.Text()
		keyword, found := searchForKeyword(allowedItems.Keywords, word)
		if found {
			keywordCounts[keyword.Keyword]++
		}
    }
	err = checkKeywordAmount(keywordCounts, allowedItems.Keywords)
	if err != nil {
		return err
	}
	return nil
}

func lintStudentCode(exercisePath string, exerciseConfig datastructures.Exercise) (err error) {
	file, err := os.Open(exercisePath)
    if err != nil {
        return fmt.Errorf("could not open %s: %w", exercisePath, err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)
    return scanStudentFile(scanner, exerciseConfig.AllowedItems)
}

func writeTemplateToFile(template, itemName string, file *os.File) error {
	content := fmt.Sprintf(template, itemName, itemName)
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error writing entry: %w", err)
	}
	return nil
}

func writeAllowedItemsLib(allowedItems datastructures.AllowedItems, file *os.File, exercise string) error {
	content := fmt.Sprintf(templates.AllowedItemsLibHeader, exercise)
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	for _, macro := range allowedItems.Macros {
		if err := writeTemplateToFile(templates.AllowedMacroTemplate, macro, file); err != nil {
			return err
		}
	}
	for _, function := range allowedItems.Functions {
		if err := writeTemplateToFile(templates.AllowedFunctionTemplate, function, file); err != nil {
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
