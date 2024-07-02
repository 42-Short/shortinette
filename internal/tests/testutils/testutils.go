package testutils

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/42-Short/shortinette/internal/errors"
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

func RunCode(executablePath string) (string, error) {
	cmd := exec.Command(executablePath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), errors.NewSubmissionError(errors.ErrRuntime, err.Error())
	}
	return stdout.String(), nil
}
