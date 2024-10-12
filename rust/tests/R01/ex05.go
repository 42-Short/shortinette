package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest05 = `
#[cfg(test)]
mod shortinette_tests_rust_0105 {
	use super::*;

	#[test]
	fn test_empty() {
		let mut v = vec![];
		deduplicate(&mut v);
		assert_eq!(v, []);
	}

	#[test]
	fn test_datatype() {
		let mut v = vec![-5i32, 0i32, 1i32, 42i32];
		deduplicate(&mut v);
		assert_eq!(v, [-5i32, 0i32, 1i32, 42i32]);
	}

	#[test]
	fn test_0() {
		let mut v = vec![1, 2, 2, 3, 4, 4];
		deduplicate(&mut v);
		assert_eq!(v, [1, 2, 3, 4]);
	}

	#[test]
	fn test_1() {
		let mut v = vec![1, 1, 1, 1];
		deduplicate(&mut v);
		assert_eq!(v, [1]);
	}

	#[test]
	fn test_2() {
		let mut v = vec![1, 2, 3, 2, 1];
		deduplicate(&mut v);
		assert_eq!(v, [1, 2, 3]);
	}

	#[test]
	fn test_3() {
		let mut v = vec![1, 2, 3, 2, 1];
		deduplicate(&mut v);
		assert_eq!(v, [1, 2, 3]);
	}

	#[test]
	fn test_4() {
		let mut v = vec![0, 3, 2, -1, -3, -5, -2, 2, -5, 2, 3, 5, 3, -5, -1, 3, 5, 0, -5, 1];
		deduplicate(&mut v);
		assert_eq!(v, [0, 3, 2, -1, -3, -5, -2, 5, 1]);
	}
}
`

func ex05Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest05, exercise.TurnInFiles[0]); err != nil {
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("05", "ex05", []string{"src/lib.rs", "Cargo.toml"}, 15, ex05Test)
}
