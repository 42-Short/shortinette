package R00

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/functioncheck"
	exercisebuilder "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
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

func compileWithRustcTestOption(dir string, turnInFile string) error {
	cmd := exec.Command("rustc", "--test", turnInFile)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Println(err)
		return errors.NewSubmissionError(errors.ErrInvalidCompilation, string(output))
	}
	logger.Info.Printf("%s/%s compiled with rustc --test\n", dir, turnInFile)
	return nil
}

func ex01Test(test *exercisebuilder.Test) bool {
	if err := functioncheck.Execute(*test, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", test.Name, err)
		return false
	}
	filePath := fmt.Sprintf("studentcode/%s/%s", test.TurnInDirectory, test.TurnInFile)
	if err := testutils.AppendStringToFile(CargoTest, filePath); err != nil {
		logger.Error.Printf("could not write to %s: %v", filePath, err)
		logger.File.Printf("internal error: could not write to %s: %v", filePath, err)
		return false
	}
	turnInDirectory := fmt.Sprintf("studentcode/%s", test.TurnInDirectory)
	if err := compileWithRustcTestOption(turnInDirectory, test.TurnInFile); err != nil {
		logger.File.Printf("[%s KO]: invalid compilation: %v", test.Name, err)
		return false
	}
	if output, err := testutils.RunCode(strings.TrimSuffix(filePath, ".rs")); err != nil {
		logger.File.Printf("[%s KO]: invalid output: %v", test.Name, output)
	}
	logger.File.Printf("[%s OK]", test.Name)
	return true
}

func ex01() exercisebuilder.ExerciseBuilder {
	return exercisebuilder.NewExerciseBuilder().
		SetName("EX01").
		SetTurnInDirectory("ex01").
		SetTurnInFile("min.rs").
		SetExerciseType("function").
		SetPrototype("min(0, 0)").
		SetAllowedMacros([]string{"println"}).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(map[string]int{"unsafe": 0}).
		SetExecuter(ex01Test)
}
