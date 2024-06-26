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
	"github.com/42-Short/shortinette/internal/functioncheck"
	"github.com/42-Short/shortinette/internal/git"
	"github.com/42-Short/shortinette/internal/logger"
)

func compileProgram(directory, turnInFile string) error {
	if _, err := os.Stat(fmt.Sprintf("%s/Cargo.toml", directory)); os.IsNotExist(err) {
		return compileWithRustc(directory, turnInFile)
	}
	return compileWithCargo(directory)
}

func compileWithRustc(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", turnInFile)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/%s compiled with rustc\n", dir, turnInFile)
	return nil
}

func compileWithCargo(dir string) error {
	cmd := exec.Command("cargo", "build")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/Cargo.toml compiled\n", dir)
	return nil
}

func compileWithRustcTestOption(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", "--test", turnInFile)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/%s compiled with rustc --test\n", dir, turnInFile)
	return nil
}

func checkAssertions(output string, assertions datastructures.Test) error {
	for _, assert := range assertions.AssertEq {
		outputReplaced := strings.ReplaceAll(output, "\n", "\\n")
		assertReplaced := strings.ReplaceAll(assert, "\n", "\\n")
		if outputReplaced != assertReplaced {
			return fmt.Errorf("assertion failed: expected '%s', got '%s'", assertReplaced, outputReplaced)
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

func prepareEnvironment(configFilePath string, repoId string, codeDirectory string) (*datastructures.Config, map[string][]datastructures.AllowedItem, error) {
	allowedItems, err := config.GetAllowedItems(configFilePath)
	if err != nil {
		return nil, nil, errors.NewInternalError(errors.ErrInternal, fmt.Sprintf("failed to get allowed items: %v", err))
	}
	conf, err := config.GetConfig(configFilePath)
	if err != nil {
		return nil, nil, errors.NewInternalError(errors.ErrInternal, fmt.Sprintf("failed to get test configuration: %v", err))
	}
	if err := git.Get(fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId), codeDirectory); err != nil {
		return nil, nil, errors.NewInternalError(errors.ErrInternal, fmt.Sprintf("failed to clone repository: %v", err))
	}
	if err := logger.InitializeTraceLogger(repoId); err != nil {
		return nil, nil, errors.NewInternalError(errors.ErrInternal, fmt.Sprintf("failed to initalize logging system: %v", err))
	}
	return conf, allowedItems, nil
}

func runProgramTests(exercise datastructures.Exercise, codeDirectory string, executablePath string) error {
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
		return errors.NewSubmissionError(errors.ErrInvalidOutput, output)
	}
	return nil
}

func runTestsForExercise(exercise datastructures.Exercise, codeDirectory string) error {
	studentCodeParentDir := fmt.Sprintf("%s/%s", codeDirectory, exercise.TurnInDirectory)
	executablePath := strings.TrimSuffix(fmt.Sprintf("%s/%s", studentCodeParentDir, exercise.TurnInFile), ".rs")

	if exercise.Type == "program" {
		if err := runProgramTests(exercise, studentCodeParentDir, executablePath); err != nil {
			return err
		}
	} else if exercise.Type == "function" {
		if err := runFunctionTests(exercise, studentCodeParentDir, executablePath); err != nil {
			return err
		}
	}
	return nil
}

func Run(configFilePath string, repoId string, codeDirectory string) (map[string]error, error) {
	defer os.RemoveAll(codeDirectory)

	conf, _, err := prepareEnvironment(configFilePath, repoId, codeDirectory)
	if err != nil {
		return nil, err
	}

	results := make(map[string]error)
	for key, exercise := range conf.Exercises {
		logger.File.Printf("[%s]\n", key)
		if err := functioncheck.Execute(exercise, repoId); err != nil {
			logger.File.Println(err)
			results[key] = err
			continue
		}
		if err := runTestsForExercise(exercise, codeDirectory); err != nil {
			logger.File.Println(err)
			results[key] = err
			continue
		}
		results[key] = nil
		logger.File.Printf("[%s] passed\n", key)
	}
	return results, nil
}
