package R04

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString04 = ``

func testMoreCommands(workingDirectory string) Exercise.Result {
	commandLine := []string{"cargo", "run", "--", "echo", "Hello", ",", "sleep", "1", ",", "touch", "World", ",", "pwd", ",", "echo", "test"}

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
		return Exercise.RuntimeError(out.err.Error(), strings.Join(commandLine, " "))
	}

	pwd, err := testutils.RunCommandLine(workingDirectory, "pwd", nil)
	if err != nil {
		return Exercise.InternalError(err.Error())
	}

	outputSlice := strings.Split(string(out.out), "\n")
	missingOutputs := map[string]bool{"Hello": true, pwd[:len(pwd)-1]: true, "test": true}
	for _, line := range outputSlice {
		delete(missingOutputs, line)
	}
	if len(missingOutputs) > 0 {
		return Exercise.AssertionError(fmt.Sprintf("*\nHello\n*\n%s\ntest\n*\n", pwd), strings.Join(outputSlice, "\n"), strings.Join(commandLine, " "))
	}

	if _, err := os.OpenFile(filepath.Join(workingDirectory, "World"), os.O_RDONLY, 0644); os.IsNotExist(err) {
		return Exercise.RuntimeError("did not create file 'World'", strings.Join(commandLine, " "))
	}

	return Exercise.Passed("OK")
}

func testTwoCommands(workingDirectory string) Exercise.Result {
	commandLine := []string{"cargo", "run", "--", "echo", "Hello", ",", "echo", "World"}

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
		return Exercise.RuntimeError(out.err.Error(), strings.Join(commandLine, " "))
	}

	outputSlice := strings.Split(string(out.out), "\n")
	missingOutputs := map[string]bool{"Hello": true, "World": true}
	for _, line := range outputSlice {
		delete(missingOutputs, line)
	}
	if len(missingOutputs) > 0 {
		return Exercise.AssertionError("*\nHello\n*\nWorld\n*\n", strings.Join(outputSlice, "\n"), commandLine...)
	}
	return Exercise.Passed("OK")
}

func testInvalidCommand(workingDirectory string) Exercise.Result {
	commandLine := []string{"cargo", "run", "--", "invalidcmd"}

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

func ex04Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString04, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testNoInput(workingDirectory); !result.Passed {
		return result
	}
	if result = testTwoCommands(workingDirectory); !result.Passed {
		return result
	}
	if result = testMoreCommands(workingDirectory); !result.Passed {
		return result
	}
	if result = testInvalidCommand(workingDirectory); !result.Passed {
		return result
	}

	return Exercise.Passed("OK")
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "ex04", []string{"Cargo.toml", "src/main.rs"}, 10, ex04Test)
}
