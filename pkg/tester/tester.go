package tester

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/42-Short/shortinette/internal/config"
	"github.com/42-Short/shortinette/internal/datastructures"
	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/functioncheck"
	"github.com/42-Short/shortinette/pkg/git"
)

func compileStudentCode(codeDir, turnInDir, turnInFile string) error {
	parentDir := fmt.Sprintf("./%s/%s/", codeDir, turnInDir)

	if _, err := os.Stat(fmt.Sprintf("%s/Cargo.toml", parentDir)); os.IsNotExist(err) {
		return compileWithRustc(parentDir, turnInFile)
	} else {
		return compileWithCargo(parentDir)
	}
}

func compileWithRustc(dir, turnInFile string) error {
	cmd := exec.Command("rustc", turnInFile)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
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

func runStudentCode(executablePath string) (string, error) {
	cmd := exec.Command(executablePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", errors.NewSubmissionError(errors.ErrRuntime, err.Error())
	}
	return out.String(), nil
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
			return nil, nil, fmt.Errorf("function check failed for %s: %w", key, err)
		}
	}

	return conf, allowedItems, nil
}

func Run(configFilePath, studentLogin, codeDirectory string) error {
	// defer os.RemoveAll(codeDirectory)

	conf, _, err := prepareEnvironment(configFilePath)
	if err != nil {
		fmt.Println(err)
	}

	if err := git.Get(fmt.Sprintf("https://github.com/42-Short/%s.git", studentLogin), codeDirectory); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	for key, exercise := range conf.Exercises {
		fmt.Printf("Running tests for %s...\n", key)

		if err := compileStudentCode(codeDirectory, exercise.TurnInDirectory, exercise.TurnInFile); err != nil {
			return err
		}

		executablePath := fmt.Sprintf("%s/%s/%s", codeDirectory, exercise.TurnInDirectory, exercise.TurnInFile)
		executablePath = strings.TrimSuffix(executablePath, ".rs")

		output, err := runStudentCode(executablePath)
		if err != nil {
			return err
		}

		if exercise.Type == "program" {
			if err := checkAssertions(output, exercise.Tests); err != nil {
				return err
			}
		} else {
			// TODO
			fmt.Println("Not implemented")
		}
		fmt.Printf("Tests for %s passed\n", key)
	}

	fmt.Println("All tests passed for all exercises!")
	return nil
}
