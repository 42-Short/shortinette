//go:build ignore
package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest01 = `
#[cfg(test)]
mod shortinette_tests_rust_0101 {
    use super::*;

    #[test]
    fn test_0() {
		let a: i32 = 1;
		let b: i32 = 2;
        assert_eq!(min(&a, &b), &a);
    }

	#[test]
    fn test_1() {
		let a: i32 = 2;
		let b: i32 = 1;
        assert_eq!(min(&a, &b), &b);
    }

	#[test]
    fn test_2() {
		let a: i32 = 1;
		let b: i32 = 1;
        assert_eq!(min(&a, &b), &a);
    }

	#[test]
    fn test_3() {
		let a: i32 = -1;
		let b: i32 = 0;
        assert_eq!(min(&a, &b), &a);
    }

	#[test]
    fn test_4() {
		let a = i32::MIN;
		let b: i32 = 1;
        assert_eq!(min(&a, &b), &a);
    }

	#[test]
    fn test_5() {
		let a = 1;
		let b = i32::MIN;
        assert_eq!(min(&a, &b), &b);
    }
}
`

func ex01Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0, "return": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest01, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"src/lib.rs", "Cargo.toml"}, 10, ex01Test)
}
