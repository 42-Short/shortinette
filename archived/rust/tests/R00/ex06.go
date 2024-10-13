//go:build ignore
package R00

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

const (
	FirstMsg = "Me and my infinite wisdom have found an appropriate secret you shall yearn for.\n"
	Equal    = "That is right! The secret was indeed the number %d, which you have brilliantly discovered!\n"
	Greater  = "This student might not be as smart as I was told. This answer is obviously too weak.\n"
	Less     = "Sometimes I wonder whether I should retire. I would have guessed higher.\n"
)

func guessingGameTest(exercise *Exercise.Exercise) (Exercise.Result, int64) {
	var number int64

	min := int64(math.MinInt32)
	max := int64(math.MaxInt32)
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cargo", "run")
	cmd.Dir = workingDirectory
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return Exercise.InternalError(fmt.Sprintf("error creating stdin pipe: %v", err.Error())), 0
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Exercise.InternalError(fmt.Sprintf("error creating stdout pipe: %v", err.Error())), 0
	}
	if err := cmd.Start(); err != nil {
		return Exercise.InternalError(fmt.Sprintf("error running command: %v", err.Error())), 0
	}
	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return Exercise.RuntimeError("Timeout"), 0
		}
		return Exercise.InternalError(fmt.Sprintf("error reading line: %v", err.Error())), 0
	}
	if line != FirstMsg {
		return Exercise.AssertionError(FirstMsg, line), 0
	}
	for {
		number = (min + max) / 2
		if _, err := writer.WriteString(fmt.Sprintf("%d\n", number)); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return Exercise.RuntimeError("Timeout"), 0
			}
			return Exercise.InternalError(fmt.Sprintf("error writing to stdin: %v", err.Error())), 0
		}
		writer.Flush()
		line, err := reader.ReadString('\n')
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return Exercise.RuntimeError("Timeout"), 0
			}
			return Exercise.InternalError(fmt.Sprintf("error reading line: %v", err.Error())), 0
		}
		if line == fmt.Sprintf(Equal, number) {
			break
		} else if line == Greater {
			if min == max {
				return Exercise.AssertionError(fmt.Sprintf(Equal, number), line), 0
			}
			min = number + 1
		} else if line == Less {
			if min == max {
				return Exercise.AssertionError(fmt.Sprintf(Equal, number), line), 0
			}
			max = number - 1
		} else {
			return Exercise.InternalError(fmt.Sprintf("unexpected output: %s", line)), 0
		}
	}
	stdin.Close()
	if err = cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return Exercise.RuntimeError("Timeout"), 0
		}
		return Exercise.InternalError(err.Error()), 0
	}
	return Exercise.Passed("OK"), number
}

func ex06Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0, "<": 0, ">": 0, "<=": 0, ">=": 0, "==": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"build"}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	result, number := guessingGameTest(exercise)
	if !result.Passed {
		return result
	}
	for i := 0; i < 3; i++ {
		result, newnumber := guessingGameTest(exercise)
		if !result.Passed {
			return result
		}
		if number != newnumber {
			return Exercise.Passed("OK")
		}
	}
	return Exercise.Result{Passed: false, Output: "numbers don't appear to be random"}
}

func ex06() Exercise.Exercise {
	return Exercise.NewExercise("06", "ex06", []string{"src/main.rs", "Cargo.toml"}, 15, ex06Test)
}
