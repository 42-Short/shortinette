package functioncheck

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/42-Short/shortinette/pkg/git"
)

type AllowedItem struct {
	Name string
	Type string
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
		return fmt.Errorf("error writing macro entry: %w", err)
	}
	return nil
}

func writeFunctionEntry(itemName string, file *os.File) error {
	content := fmt.Sprintf(`
	#[cfg(not(feature = "allowed_%s"))]
	pub fn %s() {}
`, itemName, itemName)
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("error writing function entry: %w", err)
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

func writeStudentCodeCargoToml(exercise string) error {
	path := "functioncheck/Cargo.toml"
	file, err := CreateFileWithDirs(path)
	if err != nil {
		return err
	}
	defer file.Close()
	content := fmt.Sprintf(`[package]
name = "functioncheck"
version = "0.1.0"
edition = "2021"

[dependencies]
allowedfunctions = { path = "allowedfunctions" }

[[bin]]
name = "functioncheck"
path = "src/%s/temp.rs"

[workspace]
`, exercise)
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func initCompilingEnvironment(allowedItems []AllowedItem, exercise string) error {
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

	if err := writeStudentCodeCargoToml(exercise); err != nil {
		return err
	}

	return nil
}

func prependHeadersToStudentCode(filePath string, exercise string) (err error) {

	originalFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	tempFilePath := fmt.Sprintf("functioncheck/src/%s/temp.rs", exercise)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	headers := fmt.Sprintf(`#![no_std]
#[macro_use]
extern crate allowedfunctions;
use allowedfunctions::%s::*;
`, exercise)
	if _, err := tempFile.WriteString(headers); err != nil {
		return err
	} else if originalContent, err := io.ReadAll(originalFile); err != nil {
		return err
	} else if _, err := tempFile.Write(originalContent); err != nil {
		return err
	}
	return nil
}

func compileWithDummyLib(sourceDir string) (string, error) {
	cmd := exec.Command("cargo", "build")
	cmd.Dir = sourceDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error compiling code in %s: %s\nCompiler Output:\n%s", sourceDir, err, output)
	}
	return string(output), nil
}

func setToSlice(forbiddenFunctionSet map[string]bool) []string {
	var slice []string
	for key := range forbiddenFunctionSet {
		slice = append(slice, key)
	}
	return slice
}

func parseForbiddenFunctions(compilerOutput string) ([]string, error) {
	re, err := regexp.Compile("error: cannot find (function|macro) `" + `(\w+)` + "` in this scope")
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %s", err)
	}

	matches := re.FindAllStringSubmatch(compilerOutput, -1)
	if matches == nil {
		return nil, fmt.Errorf("no forbidden functions found")
	}

	forbiddenFunctionsSet := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 2 {
			forbiddenFunctionsSet[match[2]] = true
		}
	}

	return setToSlice(forbiddenFunctionsSet), nil
}

func Execute(allowedItems []AllowedItem, exercise string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error: %w", err)
		}
	}()

	if err = initCompilingEnvironment(allowedItems, exercise); err != nil {
		return err
	} else if err = git.Execute("https://github.com/42-Short/abied-ch-R00.git", "functioncheck/src/"); err != nil {
		return err
	}

	studentCodeFilePath := fmt.Sprintf("functioncheck/src/%s/main.rs", exercise)
	err = prependHeadersToStudentCode(studentCodeFilePath, exercise)
	if err != nil {
		return fmt.Errorf("error prepending headers to student code: %s", err)
	}

	output, err := compileWithDummyLib("functioncheck/")
	if err != nil {
		usedForbiddenFunctions, parseErr := parseForbiddenFunctions(output)
		if parseErr != nil {
			return fmt.Errorf("error parsing forbidden functions: %v", parseErr)
		}
		if len(usedForbiddenFunctions) > 0 {
			forbiddenFunctionsStr := strings.Join(usedForbiddenFunctions, ", ")
			return fmt.Errorf("error: forbidden functions used: %s", forbiddenFunctionsStr)
		}
	}
	return nil
}
