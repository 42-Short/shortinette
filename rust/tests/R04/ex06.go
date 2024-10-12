package R04

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString06 = ``

var expectedOutputOptionm20 = `/lib64/ld-linux-x86-64.so.2
_ITM_deregisterTMCloneTable
_ITM_registerTMCloneTable
GCC: (Debian 12.2.0-14) 12.2.0
deregister_tm_clones
__do_global_dtors_aux
__do_global_dtors_aux_fini_array_entry
__frame_dummy_init_array_entry
_GLOBAL_OFFSET_TABLE_
__libc_start_main@GLIBC_2.34
_ITM_deregisterTMCloneTable
_ITM_registerTMCloneTable
__cxa_finalize@GLIBC_2.2.5
`

// TODO: use strings command instead of hardcoding output
func testNonExisting(workingDirectory string) Exercise.Result {
	commandLine := "cargo run do_not_panic_but_this_file_does_not_exist"
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine}); err != nil {
		if strings.Contains(err.Error(), "thread 'main' panicked") {
			return Exercise.RuntimeError(err.Error(), commandLine)
		}
	}
	return Exercise.Passed("OK")
}

func testOutput(workingDirectory string) Exercise.Result {
	if err := os.WriteFile(filepath.Join(workingDirectory, "test.c"), []byte("int main(){return 0;}"), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "cc", []string{"test.c"}); err != nil {
		return Exercise.InternalError(err.Error())
	}
	commandLine := "cargo run a.out -m 20"
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine})
	if err != nil {
		return Exercise.RuntimeError(err.Error(), commandLine)
	}
	if output != expectedOutputOptionm20 {
		return Exercise.AssertionError(expectedOutputOptionm20, output, "echo 'int main(){return 0;}' > test.c", "cc test.c", commandLine)
	}
	return Exercise.Passed("OK")
}

func testMaxSize(workingDirectory string) Exercise.Result {
	commandLine := "cargo run a.out -M 2"
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine})
	if err != nil {
		return Exercise.RuntimeError(err.Error(), commandLine)
	}
	outputLines := strings.Split(output, "\n")
	for _, line := range outputLines {
		if len(line) > 2 {
			return Exercise.AssertionError(fmt.Sprintf("line '%s' is too long", line), "max len 2", commandLine)
		}
	}
	return Exercise.Passed("OK")
}

func testMinSize(workingDirectory string) Exercise.Result {
	commandLine := "cargo run a.out -m 20"
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine})
	if err != nil {
		return Exercise.RuntimeError(err.Error(), commandLine)
	}
	outputLines := strings.Split(output, "\n")
	for _, line := range outputLines {
		if line != "" && len(line) < 20 {
			return Exercise.AssertionError(fmt.Sprintf("line '%s' is too short", line), "min len 20", commandLine)
		}
	}
	return Exercise.Passed("OK")
}

func ex06Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString06, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testNoInput(workingDirectory); !result.Passed {
		return result
	}
	// if result = testOutput(workingDirectory); !result.Passed {
	// 	return result
	// }
	if result = testNonExisting(workingDirectory); !result.Passed {
		return result
	}
	if result = testMaxSize(workingDirectory); !result.Passed {
		return result
	}
	if result = testMinSize(workingDirectory); !result.Passed {
		return result
	}

	return Exercise.Passed("OK")
}

func ex06() Exercise.Exercise {
	return Exercise.NewExercise("06", "ex06", []string{"Cargo.toml", "src/main.rs"}, 15, ex06Test)
}
