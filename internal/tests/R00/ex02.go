package R00

import (
	"time"

	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	IExercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

const YesMain = `
fn main() {
	yes();
}
`

func yesBuilder() Exercise.Exercise {
	return IExercise.NewExercise("EX02", "ex02", "yes.rs", "function", "yes()", []string{"println"}, nil, nil, nil)
}

func yes() bool {
	exercise := yesBuilder()
	fullTurnInFilePath := testutils.FullTurnInFilePath("studentcode", exercise)
	if err := testutils.AppendStringToFile(YesMain, fullTurnInFilePath); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	executablePath := testutils.ExecutablePath(fullTurnInFilePath, ".rs")
	_, err := testutils.RunCode(executablePath, testutils.WithTimeout(1*time.Second))
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	return true
}

func collatz() bool {
	return true
}

func print_bytes() bool {
	return true
}

func ex02Test(exercise *IExercise.Exercise) bool {
	if yes() && collatz() && print_bytes() {
		return true
	}
	return false
}

func ex02() Exercise.Exercise {
	return IExercise.NewExercise("EX02", "ex02", "", "", "", nil, nil, nil, ex02Test)
}
