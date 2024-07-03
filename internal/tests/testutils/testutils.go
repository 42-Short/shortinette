package testutils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	IExercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
)

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

// WithTimeout allows setting a timeout for the code execution
func WithTimeout(d time.Duration) RunCodeOption {
	return func(cmd *exec.Cmd) {
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		_ = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
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
		return stderr.String(), fmt.Errorf("runtime error: %v", err)
	}
	return stdout.String(), nil
}

func FullTurnInFilePath(codeDirectory string, exercise IExercise.Exercise) string {
	return fmt.Sprintf("%s/%s/%s", codeDirectory, exercise.TurnInDirectory, exercise.TurnInFile)
}

func ExecutablePath(fullTurnInFilePath string, suffix string) string {
	return strings.TrimSuffix(fullTurnInFilePath, suffix)
}
