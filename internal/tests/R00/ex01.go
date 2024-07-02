package R00

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/functioncheck"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/testbuilder"
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

func appendMainToFile(dest string) error {
	destFile, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = destFile.WriteString(CargoTest); err != nil {
		return err
	}
	return nil
}

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

func runCode(executablePath string) (string, error) {
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

func ex01Test(test *testbuilder.Test) bool {
	if err := functioncheck.Execute(*test, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", test.Name, err)
		return false
	}
	filePath := fmt.Sprintf("studentcode/%s/%s", test.TurnInDirectory, test.TurnInFile)
	if err := appendMainToFile(filePath); err != nil {
		logger.Error.Printf("could not write to %s: %v", filePath, err)
		logger.File.Printf("internal error: could not write to %s: %v", filePath, err)
		return false
	}
	turnInDirectory := fmt.Sprintf("studentcode/%s", test.TurnInDirectory)
	if err := compileWithRustcTestOption(turnInDirectory, test.TurnInFile); err != nil {
		logger.File.Printf("[%s KO]: invalid compilation: %v", test.Name, err)
		return false
	}
	if output, err := runCode(strings.TrimSuffix(filePath, ".rs")); err != nil {
		logger.File.Printf("[%s KO]: invalid output: %v", test.Name, output)
	}
	logger.File.Printf("[%s OK]", test.Name)
	return true
}

func ex01() testbuilder.TestBuilder {
	return testbuilder.NewTestBuilder().
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
