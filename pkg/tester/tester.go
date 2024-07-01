package tester

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/42-Short/shortinette/internal/config"
	"github.com/42-Short/shortinette/internal/datastructures"
	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
)

func compileProgram(directory, turnInFile string) error {
	if _, err := os.Stat(fmt.Sprintf("%s/Cargo.toml", directory)); os.IsNotExist(err) {
		return compileWithRustc(directory, turnInFile)
	} else {
		return compileWithCargo(directory)
	}
}

func compileWithRustc(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", turnInFile)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	return nil
}

func compileWithCargo(dir string) error {
	cmd := exec.Command("cargo", "build")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	return nil
}

func compileWithRustcTestOption(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", "--test", turnInFile)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	fmt.Println("Compiled code with test option")
	return nil
}

func checkAssertions(output string, assertions datastructures.Test) error {
	for _, assert := range assertions.AssertEq {
		if output != assert {
			return fmt.Errorf("assertion failed: expected %s, got %s", assert, output)
		}
	}
	for _, assert := range assertions.AssertNe {
		if output == assert {
			return fmt.Errorf("assertion failed: expected not %s, but got %s", assert, output)
		}
	}
	return nil
}

func runCode(executablePath string) (string, error) {
	cmd := exec.Command(executablePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.NewSubmissionError(errors.ErrRuntime, err.Error())
	}
	return out.String(), nil
}

func appendToFile(source string, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return err
	}
	return nil
}

func prepareEnvironment(configFilePath string) (*datastructures.Config, map[string][]datastructures.AllowedItem, error) {
	allowedItems, err := config.GetAllowedItems(configFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get allowed items: %w", err)
	}
	
	conf, err := config.GetConfig(configFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get config: %w", err)
	}

	for key, exercise := range conf.Exercises {
		if err := functioncheck.Execute(exercise); err != nil {
			return conf, allowedItems, fmt.Errorf("function check failed for %s: %w", key, err)
		}
	}
	return conf, allowedItems, nil
}

func runProgramTests(exercise datastructures.Exercise, codeDirectory string, executablePath string) error {
	fmt.Println(codeDirectory, exercise.TurnInDirectory, exercise.TurnInFile)
	if err := compileProgram(codeDirectory, exercise.TurnInFile); err != nil {
		return err
	}
	output, err := runCode(executablePath)
	if err != nil {
		return err
	}
	if err := checkAssertions(output, exercise.Tests); err != nil {
		return err
	}
	return nil
}

func runFunctionTests(exercise datastructures.Exercise, codeDirectory string, executablePath string) (err error) {
	if err = appendToFile(exercise.TestsPath, fmt.Sprintf("%s/min.rs", codeDirectory)); err != nil {
		return err
	}
	if err = compileWithRustcTestOption(codeDirectory, exercise.TurnInFile); err != nil {
		return err
	}
	if output, err := runCode(executablePath); err != nil {
		return errors.NewSubmissionError(errors.ErrFailedTests, output)
	}
	return nil
}

func runTestsForExercise(exercise datastructures.Exercise, codeDirectory string, exerciseNumber string) error {
	fmt.Printf("Running tests for %s...\n", exerciseNumber)

	studentCodeParentDir := fmt.Sprintf("%s/%s", codeDirectory, exercise.TurnInDirectory)
	executablePath := strings.TrimSuffix(fmt.Sprintf("%s/%s", studentCodeParentDir, exercise.TurnInFile), ".rs")

	if exercise.Type == "program" {
		if err := runProgramTests(exercise, studentCodeParentDir, executablePath); err != nil {
			fmt.Println(err)
		}
	} else if exercise.Type == "function" {
		if err := runFunctionTests(exercise, studentCodeParentDir, executablePath); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("Tests for %s passed\n", executablePath)
	return nil
}

func Run(configFilePath, studentLogin, codeDirectory string) error {
	defer os.RemoveAll(codeDirectory)

	conf, _, err := prepareEnvironment(configFilePath)
	if err != nil {
		fmt.Println(err)
	}

	if err := git.Get(fmt.Sprintf("https://github.com/%s/%s.git",os.Getenv("GITHUB_ORGANISATION"), studentLogin), codeDirectory); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	for key, exercise := range conf.Exercises {
		if err := runTestsForExercise(exercise, codeDirectory, key); err != nil {
			return err
		}
	}

	fmt.Println("All tests passed for all exercises!")
	return nil
}
