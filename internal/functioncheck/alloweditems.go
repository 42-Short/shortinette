package functioncheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/internal/templates"
	"github.com/42-Short/shortinette/internal/testbuilder"
)

func searchForKeyword(keywords map[string]int, word string) (keyword string, found bool) {
	for keyword := range keywords {
		if word == keyword {
			return keyword, true
		}
	}
	return keyword, false
}

func checkKeywordAmount(keywordCounts map[string]int, keywords map[string]int) (err error) {
	foundKeywords := make([]string, 0, len(keywords))
	for keyword, allowedAmount := range keywords {
		if count, inMap := keywordCounts[keyword]; inMap {
			if count > allowedAmount {
				foundKeywords = append(foundKeywords, keyword)
			}
		}
	}
	if len(foundKeywords) > 0 {
		return fmt.Errorf("keywords %s are used more often than allowed", strings.Join(foundKeywords, ", "))
	}
	return nil
}

func scanStudentFile(scanner *bufio.Scanner, allowedKeywords map[string]int) (err error) {
	keywordCounts := make(map[string]int)
	for scanner.Scan() {
		word := scanner.Text()
		keyword, found := searchForKeyword(allowedKeywords, word)
		if found {
			keywordCounts[keyword]++
		}
	}
	err = checkKeywordAmount(keywordCounts, allowedKeywords)
	if err != nil {
		return err
	}
	return nil
}

func lintStudentCode(exercisePath string, test testbuilder.Test) (err error) {
	file, err := os.Open(exercisePath)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", exercisePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	return scanStudentFile(scanner, test.AllowedKeywords)
}

func writeTemplateToFile(template, itemName string, file *os.File) error {
	content := fmt.Sprintf(template, itemName, itemName)
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error writing entry: %w", err)
	}
	return nil
}

func writeAllowedItemsLib(test testbuilder.Test, file *os.File, exercise string) error {
	content := fmt.Sprintf(templates.AllowedItemsLibHeader, exercise)
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	for _, macro := range test.AllowedMacros {
		if err := writeTemplateToFile(templates.AllowedMacroTemplate, macro, file); err != nil {
			return err
		}
	}
	for _, function := range test.AllowedFunctions {
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
