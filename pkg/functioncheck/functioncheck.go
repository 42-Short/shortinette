package functioncheck

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/42-Short/shortinette/pkg/git"
)

func initCompilingEnvironment(allowedItems []AllowedItem, exercise string) error {
	libFilePath := "internal/allowedfunctions/src/lib.rs"
	file, err := createFileWithDirs(libFilePath)
	if err != nil {
		return err
	}

	if err := writeAllowedItemsLib(allowedItems, file); err != nil {
		return err
	}

	if err := writeCargoToml("internal/allowedfunctions/Cargo.toml", allowedItemsCargoToml); err != nil {
		return err
	}

	cargoTomlContent := fmt.Sprintf(cargoTomlTemplate, "internal", exercise)
	if err := writeCargoToml("internal/Cargo.toml", cargoTomlContent); err != nil {
		return err
	}

	return nil
}

func prependHeadersToStudentCode(filePath, exercise string) error {
	originalFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open original file: %w", err)
	}
	defer originalFile.Close()

	tempFilePath := fmt.Sprintf("internal/src/%s/temp.rs", exercise)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	defer tempFile.Close()

	headers := fmt.Sprintf(studentCodePrefix, exercise)
	if _, err := tempFile.WriteString(headers); err != nil {
		return fmt.Errorf("could not write headers: %w", err)
	}

	originalContent, err := io.ReadAll(originalFile)
	if err != nil {
		return fmt.Errorf("could not read original file content: %w", err)
	}
	if _, err := tempFile.Write(originalContent); err != nil {
		return fmt.Errorf("could not write original content to temp file: %w", err)
	}
	return nil
}

func compileWithDummyLib(sourceDir string) (string, error) {
	cmd := exec.Command("cargo", "build")
	cmd.Dir = sourceDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error compiling code in %s: %w\nCompiler Output:\n%s", sourceDir, err, output)
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
		return nil, fmt.Errorf("error compiling regex: %w", err)
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

func handleCompileError(output string) error {
	usedForbiddenFunctions, parseErr := parseForbiddenFunctions(output)
	if parseErr != nil {
		return fmt.Errorf("could not parse forbidden functions: %w", parseErr)
	} else if len(usedForbiddenFunctions) > 0 {
		forbiddenFunctions := strings.Join(usedForbiddenFunctions, ", ")
		return fmt.Errorf("forbidden functions used: %s", forbiddenFunctions)
	} else {
		return fmt.Errorf("could not compile code: %s", output)
	}
}

func Execute(allowedItems []AllowedItem, exercise string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error: %w", err)
		}
	}()

	if err = initCompilingEnvironment(allowedItems, exercise); err != nil {
		return err
	}

	if err = git.Execute("https://github.com/42-Short/abied-ch-R00.git", "internal/src/"); err != nil {
		return err
	}

	err = prependHeadersToStudentCode(fmt.Sprintf("internal/src/%s/main.rs", exercise), exercise)
	if err != nil {
		return err
	}

	output, compileErr := compileWithDummyLib("internal/")
	if compileErr != nil {
		return handleCompileError(output)
	}
	return nil
}
