package R00

import (
	"path/filepath"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var TestMod = `
#[cfg(test)]
mod test {
	use super::*;

	#[test]
	fn test1() {
		assert!(strpcmp(b"abc", b"abc"));
	}

	#[test]
	fn test2() {
		assert!(strpcmp(b"abcd", b"ab*"));
	}
	
	#[test]
	fn test3() {
		assert!(!strpcmp(b"cab", b"ab*"));
	}
	
	#[test]
	fn test4() {
		assert!(strpcmp(b"dcab", b"*ab"));
	}
	
	#[test]
	fn test5() {
		assert!(!strpcmp(b"abc", b"*ab"));
	}

	#[test]
	fn test6() {
		assert!(strpcmp(b"ab000cd", b"ab*cd"));
	}
	
	#[test]
	fn test7() {
		assert!(strpcmp(b"abcd", b"ab*cd"));
	}

	#[test]
	fn test8() {
		assert!(strpcmp(b"", b"****"));
	}
}
`

func ex07Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)

	if err := testutils.AppendStringToFile(TestMod, exercise.TurnInFiles[1]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test"})
	if err != nil {
		return Exercise.AssertionError("", output)
	}
	return Exercise.Passed("OK")
}

func ex07() Exercise.Exercise {
	return Exercise.NewExercise("07", "studentcode", "ex07", []string{"src/main.rs", "src/lib.rs", "Cargo.toml"}, "", "", []string{"assert", "assert_eq"}, nil, nil, 25, ex07Test)
}
