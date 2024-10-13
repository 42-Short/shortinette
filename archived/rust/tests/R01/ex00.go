package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest00 = `
#[cfg(test)]
mod shortinette_tests_rust_0100 {
    use super::*;

    #[test]
    fn test_add_0() {
		let a: i32 = 21;
        assert_eq!(add(&a, 21i32), 42i32);
    }

	#[test]
    fn test_add_1() {
		let a: i32 = i32::MAX;
        assert_eq!(add(&a, i32::MIN), -1i32);
    }

    #[test]
    fn test_add_assign_0() {
        let mut a: i32 = 21;
        add_assign(&mut a, 21i32);
        assert_eq!(a, 42);
    }

	#[test]
    fn test_add_assign_1() {
        let mut a: i32 = i32::MAX;
        add_assign(&mut a, i32::MIN);
        assert_eq!(a, -1);
    }
}
`

func ex00Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest00, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "ex00", []string{"src/lib.rs", "Cargo.toml"}, 10, ex00Test)
}
