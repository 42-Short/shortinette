package functioncheck

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/pkg/git"
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

func writeMacroEntry(itemName string, file *os.File) error {
	content := fmt.Sprintf(`
	#[cfg(not(feature = "allowed_%s"))]
	#[macro_export]
	macro_rules! %s {
		($($arg:tt)*) => {{}}
	}
`, itemName, itemName)
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func writeFunctionEntry(itemName string, file *os.File) error {
	content := fmt.Sprintf(`
	#[cfg(not(feature = "allowed_%s"))]
	pub fn %s() {}
`, itemName, itemName)
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func writeAllowedItemsLib(allowedItems []AllowedItem, file *os.File) error {
	exerciseNumber := "00"
	content := fmt.Sprintf("pub mod ex%s { ", exerciseNumber)
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	for _, item := range allowedItems {
		if item.Type == "macro" {
			if err := writeMacroEntry(item.Name, file); err != nil {
				return err
			}
		} else if item.Type == "function" {
			if err := writeFunctionEntry(item.Name, file); err != nil {
				return err
			}
		}
	}
	file.WriteString("}")
	return nil
}

func CreateFileWithDirs(filePath string) (*os.File, error) {
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

func writeAllowedItemsLibCargoToml() error {
	path := "functioncheck/allowedfunctions/Cargo.toml"
	file, err := CreateFileWithDirs(path)
	if err != nil {
		return err
	}
	defer file.Close()
	content := `[package]
name = "allowedfunctions"
version = "0.1.0"
edition = "2021"
`
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func writeStudentCodeCargoToml() error {
	path := "functioncheck/Cargo.toml"
	file, err := CreateFileWithDirs(path)
	if err != nil {
		return err
	}
	defer file.Close()
	content := `[package]
name = "functioncheck"
version = "0.1.0"
edition = "2021"

[dependencies]
allowedfunctions = { path = "allowedfunctions" }

[[bin]]
name = "functioncheck"
path = "src/ex00/main.rs"

[workspace]
`
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func initCompilingEnvironment(allowedItems []AllowedItem) error {
	libFilePath := "functioncheck/allowedfunctions/src/lib.rs"
	file, err := CreateFileWithDirs(libFilePath)
	if err != nil {
		return err
	}

	if err := writeAllowedItemsLib(allowedItems, file); err != nil {
		return err
	}

	if err := writeAllowedItemsLibCargoToml(); err != nil {
		return err
	}

	if err := writeStudentCodeCargoToml(); err != nil {
		return err
	}

	return nil
}

func prependHeadersToStudentCode(filePath string) error {
	originalFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	tempFilePath := "functioncheck/src/temp.rs"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	headers := `#![no_std]
#[macro_use]
extern crate allowedfunctions;
use allowedfunctions::ex00::*;
`
	if _, err := tempFile.WriteString(headers); err != nil {
		return err
	}

	originalContent, err := io.ReadAll(originalFile)
	if err != nil {
		return err
	}
	if _, err := tempFile.Write(originalContent); err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, filePath); err != nil {
		return err
	}

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
	err = git.Execute("https://github.com/42-Short/abied-ch-R00.git", "functioncheck/src/")
	if err != nil {
		return fmt.Errorf("error executing git: %s", err)
	}

	studentCodeFilePath := "functioncheck/src/ex00/main.rs"
	err = prependHeadersToStudentCode(studentCodeFilePath)
	if err != nil {
		return fmt.Errorf("error prepending headers to student code: %s", err)
	}
	return nil
}
