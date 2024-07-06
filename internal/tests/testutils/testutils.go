package testutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/42-Short/shortinette/internal/functioncheck"
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
)

func CompileWithRustc(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", turnInFile)
	cmd.Dir = dir

	_, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("invalid compilation: %s", err)
	}
	logger.Info.Printf("%s/%s compiled with rustc\n", dir, turnInFile)
	return nil
}

func ForbiddenItemsCheck(exercise Exercise.Exercise, repoId string) error {
	if err := functioncheck.Execute(exercise, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return err
	}
	return nil
}

func AssertionErrorString(testName string, expected string, got string) string {
	expectedReplaced := strings.ReplaceAll(expected, "\n", "\\n")
	gotReplaced := strings.ReplaceAll(got, "\n", "\\n")
	outputComparison := fmt.Sprintf("invalid output: expected '%s', got '%s'", expectedReplaced, gotReplaced)
	return fmt.Sprintf("[%s KO]: %v", testName, outputComparison)
}

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

type RunCodeOption func(*exec.Cmd)

var ErrTimeout = errors.New("command timed out")

// WithTimeout allows setting a timeout for the code execution
func WithTimeout(d time.Duration) RunCodeOption {
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
func RunCode(executablePath string, options ...RunCodeOption) (string, error) {
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

func containsString(hayStack []string, needle string) bool {
	for _, str := range hayStack {
		if str == needle {
			return true
		}
	}
	return false
}

func TurnInFilesCheck(exercise Exercise.Exercise) {
	fullTurnInFilesPaths := FullTurnInFilesPath(exercise)
	parentDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	err := filepath.Walk(parentDirectory, func(path string, info os.FileInfo, err error) error {
		if filepath.Base(path)[0] == '.' || path == parentDirectory || info.IsDir() {
			return nil
		} else if !containsString(fullTurnInFilesPaths, path) {
			return fmt.Errorf("'%s' not in allowed turn in files", path)
		}
		return nil
	})
	if err != nil {
		logger.Error.Printf("walk error: %v", err)
	}
}

func FullTurnInFilesPath(exercise Exercise.Exercise) []string {
	var fullFilePaths []string

	for _, path := range exercise.TurnInFiles {
		fullPath := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory, path)
		fullFilePaths = append(fullFilePaths, fullPath)
	}
	return fullFilePaths
}

func FullTurnInFilePath(codeDirectory string, exercise Exercise.Exercise, turnInFile string) string {
	return fmt.Sprintf("%s/%s/%s", codeDirectory, exercise.TurnInDirectory, turnInFile)
}

func FullTurnInDirectory(codeDirectory string, exercise Exercise.Exercise) string {
	return fmt.Sprintf("%s/%s", codeDirectory, exercise.TurnInDirectory)
}

func ExecutablePath(fullTurnInFilePath string, suffix string) string {
	return strings.TrimSuffix(fullTurnInFilePath, suffix)
}
