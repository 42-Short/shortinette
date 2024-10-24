//go:build ignore
package R04

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString05 = ``

func testValidAddress(workingDirectory string) Exercise.Result {
	commandLine := []string{"cargo", "run", "https://postman-echo.com/get"}

	if _, err := testutils.RunCommandLine(workingDirectory, commandLine[0], commandLine[1:]); err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func testUnresolvableAddress(workingDirectory string) Exercise.Result {
	commandLine := []string{"cargo", "run", "you.shall.not.panic"}

	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan outputChannel)
	go func() {
		defer wg.Done()
		cmd := exec.Command(commandLine[0], commandLine[1:]...)
		cmd.Dir = workingDirectory
		out, err := cmd.CombinedOutput()
		ch <- outputChannel{out, err}
	}()
	out := <-ch
	wg.Wait()

	if out.err != nil {
		if strings.Contains(string(out.out), "thread 'main' panicked") {
			return Exercise.RuntimeError(fmt.Sprintf("i said don't panic :(\n%s", out.err.Error()), strings.Join(commandLine, " "))
		}
	}
	return Exercise.Passed("OK")
}

func ex05Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString05, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testNoInput(workingDirectory); !result.Passed {
		return result
	}
	if result = testUnresolvableAddress(workingDirectory); !result.Passed {
		return result
	}
	if result = testValidAddress(workingDirectory); !result.Passed {
		return result
	}
	return Exercise.Passed("OK")
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("05", "ex05", []string{"Cargo.toml", "src/main.rs"}, 15, ex05Test)
}
