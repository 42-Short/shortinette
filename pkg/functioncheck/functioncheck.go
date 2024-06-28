package functioncheck

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/42-Short/shortinette/internal/datastructures"
	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/git"
)

func initCompilingEnvironment(allowedItems datastructures.AllowedItems, exercise string) error {
	libFilePath := "compile-environment/allowedfunctions/src/lib.rs"
	file, err := createFileWithDirs(libFilePath)
	if err != nil {
		return err
	}

	if err := writeAllowedItemsLib(allowedItems, file, exercise); err != nil {
		return err
	}

	if err := writeCargoToml("compile-environment/allowedfunctions/Cargo.toml", allowedItemsCargoToml); err != nil {
		return err
	}

	cargoTomlContent := fmt.Sprintf(cargoTomlTemplate, "compile-environment", exercise)
	if err := writeCargoToml("compile-environment/Cargo.toml", cargoTomlContent); err != nil {
		return err
	}

	return nil
}

func prependHeadersToStudentCode(filePath, exerciseNumber string, exerciseType string, dummyCall string) error {
	originalFile, err := os.Open(filePath)
	fmt.Printf("File Path: %s\n", filePath)
	if err != nil {
		return fmt.Errorf("could not open original file: %w", err)
	}
	defer originalFile.Close()

	tempFilePath := fmt.Sprintf("compile-environment/src/%s/temp.rs", exerciseNumber)
	fmt.Printf("Temp File Path: %s\n", tempFilePath)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	defer tempFile.Close()

	headers := fmt.Sprintf(studentCodePrefix, exerciseNumber)
	fmt.Printf("Headers: %s\n", headers)

	if _, err := tempFile.WriteString(headers); err != nil {
		return fmt.Errorf("could not write headers: %w", err)
	}
	originalContent, err := io.ReadAll(originalFile)
	fmt.Printf("Original COntent: %s\n", originalContent)
	if err != nil {
		return fmt.Errorf("could not read original file content: %w", err)
	}
	if _, err := tempFile.Write(originalContent); err != nil {
		return fmt.Errorf("could not write original content to temp file: %w", err)
	}
	if exerciseType == "function" {
		main := fmt.Sprintf(dummyMain, dummyCall)
		fmt.Printf("Main: %s\n", main)

		if _, err := tempFile.Write([]byte(main)); err != nil {
			return fmt.Errorf("could not write dummy main to temp file: %w", err)
		}
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
	re, err := regexp.Compile(`error: cannot find (function|macro) ` + `(\w+)` + ` in this scope`)
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
		return errors.NewSubmissionError(errors.ErrForbiddenItem, forbiddenFunctions)
	} else {
		return fmt.Errorf("could not compile code: %s", output)
	}
}

func Execute(exerciseConfig datastructures.Exercise) (err error) {
	// defer func() {
	// 	rmErr := os.RemoveAll("compile-environment/")
	// 	if rmErr != nil {
	// 		err = fmt.Errorf("failed to remove compile environment: %w", rmErr)
	// 	}
	// }()

	if err = initCompilingEnvironment(exerciseConfig.AllowedItems, exerciseConfig.TurnInDirectory); err != nil {
		return err
	}

	if err = git.Get("https://github.com/42-Short/shortinette-test.git", "compile-environment/src/"); err != nil {
		return err
	}

	exercisePath := fmt.Sprintf("compile-environment/src/%s/%s", exerciseConfig.TurnInDirectory, exerciseConfig.TurnInFile)
	err = prependHeadersToStudentCode(exercisePath, exerciseConfig.TurnInDirectory, exerciseConfig.Type, exerciseConfig.DummyCall)
	if err != nil {
		return err
	}

	output, compileErr := compileWithDummyLib("compile-environment/")
	if compileErr != nil {
		return handleCompileError(output)
	}

	fmt.Println("No forbidden items/keywords found")

	return nil
}
