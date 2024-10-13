package R05

import (
	"os"
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString = `
disallowed-macros = ["std::vec"]
disallowed-methods = ["std::iter::Iterator::collect", "std::iter::repeat", "std::collections::VecDeque", "std::collections::LinkedList", "std::collections::has_map::HashMap"]
disallowed-types = ["std::vec::Vec", "std::iter::Iterator", "std::collections::VecDeque", "std::collections::LinkedList", "std::collections::has_map::HashMap", "std::collections::hash_set::HashSet", "std::collections::BTreeSet", "std::collections::BinaryHeap"]
`

func ex03Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, clippyTomlAsString, map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	Ex03TestMod, err := os.ReadFile("internal/tests/R05/ex03.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := testutils.AppendStringToFile(string(Ex03TestMod), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module05_ex03_0001", "--", workingDirectory + "/target/release/ex03"}, testutils.WithTimeout(time.Minute*2))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "ex03", []string{"src/main.rs", "Cargo.toml"}, 10, ex03Test)
}
