package R00

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/functioncheck"
	"github.com/42-Short/shortinette/internal/logger"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

const CargoTest = `
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_0() {
        assert_eq!(min(1, 2), 1);
    }

    #[test]
    fn test_1() {
        assert_eq!(min(2, 1), 1);
    }

    #[test]
    fn test_2() {
        assert_eq!(min(1, 1), 1);
    }

    #[test]
    fn test_3() {
        assert_eq!(min(-1, 0), -1);
    }
}
`

func compileWithRustcTestOption(turnInFile string) error {
	cmd := exec.Command("rustc", "--test", filepath.Base(turnInFile))
	cmd.Dir = filepath.Dir(turnInFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/%s compiled with rustc --test\n", cmd.Dir, turnInFile)
	return nil
}

func ex01Test(exercise *Exercise.Exercise) bool {
	if !testutils.TurnInFilesCheck(*exercise) {
		return false
	}
	if err := functioncheck.Execute(*exercise, "shortinette-test-R00"); err != nil {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)
	if err := testutils.AppendStringToFile(CargoTest, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		logger.File.Printf("internal error: could not write to %s: %v", exercise.TurnInFiles[0], err)
		return false
	}
	if err := compileWithRustcTestOption(exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s KO]: invalid compilation: %v", exercise.Name, err)
		return false
	}
	if output, err := testutils.RunExecutable(strings.TrimSuffix(exercise.TurnInFiles[0], ".rs")); err != nil {
		logger.File.Printf("[%s KO]: invalid output: %v", exercise.Name, output)
	}
	return true
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("EX01", "studentcode", "ex01", []string{"min.rs"}, "function", "min(0, 0)", []string{"println"}, nil, map[string]int{"unsafe": 0}, ex01Test)
}
