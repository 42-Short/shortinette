// Package testutils provides utility functions for compiling, running, and managing
// code submissions, particularly for Rust, and interacting with the command line.
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

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
)

// CompileWithRustc compiles a Rust file using the rustc compiler.
//
//   - turnInFile: the path to the Rust file to be compiled
//
// Returns an error if the compilation fails.
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

// AppendStringToFile appends a source string to a destination file.
//
//   - source: the string to append to the file
//   - destFilePath: the path to the file to which the string will be appended
//
// Returns an error if the file cannot be opened or the string cannot be written.
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

// FullTurnInFilesPath constructs the full file paths for all files in an exercise's TurnInFiles.
//
//   - exercise: the Exercise struct containing the necessary directory information
//
// Returns a slice of strings representing the full file paths.
func FullTurnInFilesPath(exercise Exercise.Exercise) []string {
	var fullFilePaths []string

	for _, path := range exercise.TurnInFiles {
		fullPath := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory, path)
		fullFilePaths = append(fullFilePaths, fullPath)
	}
	return fullFilePaths
}

// DeleteStringFromFile deletes the first occurrence of a target string from a specified file.
//
//   - targetString: the string to delete from the file
//   - filePath: the path to the file from which the string will be deleted
//
// Returns an error if the file cannot be read or written.
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

// RunExecutableOption is a type for options that modify the behavior of RunExecutable.
type RunExecutableOption func(*exec.Cmd)

// ErrTimeout is an error that indicates a command has timed out.
var ErrTimeout = errors.New("command timed out")

// WithRealTimeOutput is a RunExecutableOption that allows the command to show output in real-time.
func WithRealTimeOutput() RunExecutableOption {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = io.MultiWriter(os.Stdout, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, cmd.Stderr)
	}
}

// WithTimeout is a RunExecutableOption that sets a timeout for the code execution.
//
//   - d: the duration before the command times out
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

// RunExecutable runs an executable file at the given path with the provided options.
//
//   - executablePath: the path to the executable file
//   - options: a variadic list of RunExecutableOptions that modify the command behavior
//
// Returns the stdout output as a string and an error if the execution fails.
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
			logger.Error.Printf("%s timed out: %v", executablePath, ctxErr)
			return stdout.String(), ErrTimeout
		}
		output := fmt.Sprintf("error executing %s: %v\nstdout: %s\nstderr: %s", executablePath, err, stdout.String(), stderr.String())
		return output, fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

// RunCommandLine runs a command line command with the provided options.
//
//   - workingDirectory: the directory in which to run the command
//   - command: the command to run
//   - args: the arguments for the command
//   - options: a variadic list of RunExecutableOptions that modify the command behavior
//
// Returns the stdout output as a string and an error if the command execution fails.
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
		output := fmt.Sprintf("stderr: %s\nstdout: %s", stdout.String(), stderr.String())
		return "", errors.New(output)
	}
	return stdout.String(), nil
}

// FullTurnInDirectory constructs the full path to the TurnInDirectory of an exercise.
//
//   - codeDirectory: the root directory for code submissions
//   - exercise: the Exercise struct containing the necessary directory information
//
// Returns a string representing the full path to the TurnInDirectory.
func FullTurnInDirectory(codeDirectory string, exercise Exercise.Exercise) string {
	return fmt.Sprintf("%s/%s", codeDirectory, exercise.TurnInDirectory)
}

// ExecutablePath constructs the path to an executable by removing a specified suffix from the full file path.
//
//   - fullTurnInFilePath: the full file path to the turned-in file
//   - suffix: the suffix to remove (e.g., ".rs" for Rust files)
//
// Returns a string representing the path to the executable.
func ExecutablePath(fullTurnInFilePath string, suffix string) string {
	return strings.TrimSuffix(fullTurnInFilePath, suffix)
}
