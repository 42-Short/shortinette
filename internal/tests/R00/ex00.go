package R00

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/functioncheck"
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

func ex00Compile(test *Exercise.Exercise) error {
	cmd := exec.Command("rustc", test.TurnInFile)
	cmd.Dir = fmt.Sprintf("studentcode/%s/", test.TurnInDirectory)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/%s compiled with rustc\n", cmd.Dir, test.TurnInFile)
	return nil
}

func runExecutable(executablePath string) (string, error) {
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

func ex00Test(test *Exercise.Exercise) bool {
	if err := functioncheck.Execute(*test, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", test.Name, err)
		return false
	}
	if err := ex00Compile(test); err != nil {
		logger.File.Printf("[%s KO]: %v", test.Name, err)
		return false
	}
	relativeFilePath := fmt.Sprintf("studentcode/%s/%s", test.TurnInDirectory, test.TurnInFile)
	executablePath := strings.TrimSuffix(relativeFilePath, ".rs")
	output, err := runExecutable(executablePath)
	if err != nil {
		logger.File.Printf("[%s KO]: %v", test.Name, err)
		return false
	}
	if output != "Hello, World!\n" {
		logger.File.Printf(testutils.AssertionErrorString(test.Name, "Hello, World\n", output))
		return false
	}
	logger.File.Printf("[%s OK]", test.Name)
	return true
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("EX00", "ex00", "hello.rs", "program", "", []string{"println"}, nil, map[string]int{"unsafe": 0}, ex00Test)
}
