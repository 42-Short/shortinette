package testutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	toml "github.com/42-Short/shortinette/internal/datastructures"
	"github.com/42-Short/shortinette/internal/logger"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func CheckCargoTomlContent(exercise Exercise.Exercise, expectedContent map[string]string) Exercise.Result {
	tomlPath := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory, "Cargo.toml")
	fieldMap, err := toml.ReadToml(tomlPath)
	if err != nil {
		logger.Error.Printf("internal error: %s", err)
		return Exercise.Result{Passed: false, Output: "internal error"}
	}
	var result = Exercise.Result{Passed: true, Output: ""}
	for key, expectedValue := range expectedContent {
		value, ok := fieldMap[key]
		if !ok {
			result.Passed = false
			result.Output = result.Output + fmt.Sprintf("\n'%s' not found in Cargo.toml", key)
		} else if value != expectedValue {
			result.Passed = false
			result.Output = result.Output + fmt.Sprintf("\nCargo.toml content mismatch, expected '%s', got '%s'", expectedValue, value)
		}
	}
	return result
}

func CompileWithRustc(turnInFile string) error {
	cmd := exec.Command("rustc", filepath.Base(turnInFile))
	cmd.Dir = filepath.Dir(turnInFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		logger.Error.Println(string(output))
		return err
	}
	logger.Info.Printf("%s/%s compiled with rustc\n", cmd.Dir, turnInFile)
	return nil
}

// Append source to destFilePath (e.g., a main for testing single funtions)
func AppendStringToFile(source string, destFilePath string) error {
	destFile, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = destFile.WriteString(source); err != nil {
		return err
	}
	return nil
}

func FullTurnInFilesPath(exercise Exercise.Exercise) []string {
	var fullFilePaths []string

	for _, path := range exercise.TurnInFiles {
		fullPath := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory, path)
		fullFilePaths = append(fullFilePaths, fullPath)
	}
	return fullFilePaths
}

// Delete the first occurrence of targetString from filePath
func DeleteStringFromFile(targetString, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	modifiedContent := strings.ReplaceAll(string(content), targetString, "")

	err = os.WriteFile(filePath, []byte(modifiedContent), 0666)
	if err != nil {
		return err
	}

	return nil
}

type RunExecutableOption func(*exec.Cmd)

var ErrTimeout = errors.New("command timed out")

// WithRealTimeOutput allows the command to show output in real-time.
func WithRealTimeOutput() RunExecutableOption {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = io.MultiWriter(os.Stdout, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, cmd.Stderr)
	}
}

// WithTimeout allows setting a timeout for the code execution
func WithTimeout(d time.Duration) RunExecutableOption {
	return func(cmd *exec.Cmd) {
		ctx, cancel := context.WithTimeout(context.Background(), d)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded {
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
					return
				}
			}
			cancel()
		}()
	}
}

// RunCode runs the executable at the given path with the provided options.
func RunExecutable(executablePath string, options ...RunExecutableOption) (string, error) {
	cmd := exec.Command(executablePath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	for _, opt := range options {
		opt(cmd)
	}

	if err := cmd.Run(); err != nil {
		if ctxErr := cmd.ProcessState.ExitCode(); ctxErr == -1 {
			return stdout.String(), ErrTimeout
		}
		return stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

// RunCommand runs the command line with the provided options.
func RunCommandLine(workingDirectory string, command string, args []string, options ...RunExecutableOption) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDirectory
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	for _, opt := range options {
		opt(cmd)
	}

	if err := cmd.Run(); err != nil {
		if ctxErr := cmd.ProcessState.ExitCode(); ctxErr == -1 {
			return stdout.String(), ErrTimeout
		}
		return stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

func FullTurnInDirectory(codeDirectory string, exercise Exercise.Exercise) string {
	return fmt.Sprintf("%s/%s", codeDirectory, exercise.TurnInDirectory)
}

func ExecutablePath(fullTurnInFilePath string, suffix string) string {
	return strings.TrimSuffix(fullTurnInFilePath, suffix)
}
