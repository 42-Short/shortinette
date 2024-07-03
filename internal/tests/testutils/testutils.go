package testutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	IExercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
)

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

func FullTurnInFilePath(codeDirectory string, exercise IExercise.Exercise) string {
	return fmt.Sprintf("%s/%s/%s", codeDirectory, exercise.TurnInDirectory, exercise.TurnInFile)
}

func FullTurnInDirectory(codeDirectory string, exercise IExercise.Exercise) string {
	return fmt.Sprintf("%s/%s", codeDirectory, exercise.TurnInDirectory)
}

func ExecutablePath(fullTurnInFilePath string, suffix string) string {
	return strings.TrimSuffix(fullTurnInFilePath, suffix)
}
