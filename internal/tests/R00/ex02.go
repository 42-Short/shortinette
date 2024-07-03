package R00

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	SubmissionError "github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/functioncheck"
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

const YesMain = `
fn main() {
	yes();
}
`

func compileWithRustc(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", turnInFile)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return SubmissionError.NewSubmissionError(SubmissionError.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/%s compiled with rustc\n", dir, turnInFile)
	return nil
}

func yesBuilder() Exercise.Exercise {
	return Exercise.NewExercise("EX02", "ex02", "yes.rs", "function", "yes()", []string{"println"}, nil, nil, nil)
}

func yes() bool {
	exercise := yesBuilder()
	if err := functioncheck.Execute(exercise, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	fullTurnInFilePath := testutils.FullTurnInFilePath("studentcode", exercise)
	if err := testutils.AppendStringToFile(YesMain, fullTurnInFilePath); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	directory := testutils.FullTurnInDirectory("studentcode", exercise)
	if err := compileWithRustc(directory, exercise.TurnInFile); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := testutils.ExecutablePath(fullTurnInFilePath, ".rs")
	output, err := testutils.RunCode(executablePath, testutils.WithTimeout(1*time.Second))
	if err != nil {
		if !errors.Is(err, testutils.ErrTimeout) {
			logger.File.Printf("[%s KO]: runtime error: %v", exercise.Name, err)
			return false
		}
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "y" {
			assertionError := testutils.AssertionErrorString(exercise.Name, "y", line)
			logger.File.Printf("[%s KO]: invalid output: %v", exercise.Name, assertionError)
			return false
		}
	}
	logger.File.Printf("[%s OK]", exercise.Name)
	return true
}

func collatz() bool {
	return true
}

func print_bytes() bool {
	return true
}

func ex02Test(exercise *Exercise.Exercise) bool {
	if yes() && collatz() && print_bytes() {
		return true
	}
	return false
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("EX02", "ex02", "", "", "", nil, nil, nil, ex02Test)
}
