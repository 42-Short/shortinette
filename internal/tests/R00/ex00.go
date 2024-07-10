package R00

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/functioncheck"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

func ex00Compile(exercise *Exercise.Exercise) error {
	cmd := exec.Command("rustc", filepath.Base(exercise.TurnInFiles[0]))
	cmd.Dir = filepath.Dir(exercise.TurnInFiles[0])

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s compiled with rustc\n", exercise.TurnInFiles[0])
	return nil
}

func runExecutable(executablePath string) (string, error) {
	cmd := exec.Command("./" + filepath.Base(executablePath))
	cmd.Dir = filepath.Dir(executablePath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), errors.NewSubmissionError(errors.ErrRuntime, err.Error())
	}
	return stdout.String(), nil
}

func ex00Test(exercise *Exercise.Exercise) bool {
	if err := functioncheck.Execute(*exercise, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)

	if err := ex00Compile(exercise); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := strings.TrimSuffix(exercise.TurnInFiles[0], filepath.Ext(exercise.TurnInFiles[0]))
	output, err := runExecutable(executablePath)
	if err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		logger.Error.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	if output != "Hello, World!\n" {
		logger.File.Printf(testutils.AssertionErrorString(exercise.Name, "Hello, World\n", output))
		return false
	}
	return true
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("EX00", "studentcode", "ex00", []string{"hello.rs"}, "program", "", []string{"println"}, nil, map[string]int{"unsafe": 0}, ex00Test)
}
