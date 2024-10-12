package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest07 = `
#[cfg(test)]
mod shortinette_tests_rust_0107 {
	use super::*;

	#[test]
	fn test_empty() {
		let mut tasks = vec![];
		assert_eq!(time_manager(&mut tasks), 0);
	}

	#[test]
	fn test_datatype() {
		let mut tasks = vec![
			Task{start_time: 1u32, end_time: 2u32, cookies: 5u32},
			Task{start_time: 3u32, end_time: 5u32, cookies: 10u32}
		];
		assert_eq!(time_manager(&mut tasks), 15u32);
	}

	#[test]
	fn test_0() {
		let mut tasks = vec![
			Task{start_time: 0, end_time: 3, cookies: 10},
			Task{start_time: 4, end_time: 5, cookies: 5},
			Task{start_time: 6, end_time: 10, cookies: 25}
		];
		assert_eq!(time_manager(&mut tasks), 40);
	}

	#[test]
	fn test_1() {
		let mut tasks = vec![
			Task{start_time: 0, end_time: 3, cookies: 10},
			Task{start_time: 3, end_time: 5, cookies: 5},
			Task{start_time: 5, end_time: 10, cookies: 25}
		];
		assert_eq!(time_manager(&mut tasks), 40);
	}

	#[test]
	fn test_2() {
		let mut tasks = vec![
			Task{start_time: 0, end_time: 5, cookies: 10},
			Task{start_time: 3, end_time: 7, cookies: 5},
			Task{start_time: 5, end_time: 10, cookies: 25}
		];
		assert_eq!(time_manager(&mut tasks), 35);
	}

	#[test]
	fn test_3() {
		let mut tasks = vec![
			Task{start_time: 0, end_time: 5, cookies: 1},
			Task{start_time: 3, end_time: 7, cookies: 30},
			Task{start_time: 5, end_time: 10, cookies: 25}
		];
		assert_eq!(time_manager(&mut tasks), 30);
	}

	#[test]
	fn test_4() {
		let mut tasks = vec![
			Task{start_time: 0, end_time: 5, cookies: 1},
			Task{start_time: 3, end_time: 7, cookies: 24},
			Task{start_time: 5, end_time: 10, cookies: 25}
		];
		assert_eq!(time_manager(&mut tasks), 26);
	}

	#[test]
	fn test_5() {
		let mut tasks = vec![
			Task{start_time: 0, end_time: 5, cookies: 1},
			Task{start_time: 3, end_time: 7, cookies: 25},
			Task{start_time: 5, end_time: 10, cookies: 25}
		];
		assert_eq!(time_manager(&mut tasks), 26);
	}

	#[test]
	fn test_6() {
		let mut tasks = vec![
			Task{start_time: 5, end_time: 10, cookies: 25},
			Task{start_time: 3, end_time: 5, cookies: 5},
			Task{start_time: 0, end_time: 3, cookies: 10}
		];
		assert_eq!(time_manager(&mut tasks), 40);
	}

	#[test]
	fn test_7() {
		let mut tasks = vec![
			Task{start_time: 2, end_time: 25, cookies: 10},
			Task{start_time: 1, end_time: 22, cookies: 23},
			Task{start_time: 6, end_time: 19, cookies: 22},
			Task{start_time: 6, end_time: 16, cookies: 24},
			Task{start_time: 2, end_time: 7, cookies: 12},
			Task{start_time: 19, end_time: 20, cookies: 20},
			Task{start_time: 16, end_time: 18, cookies: 23},
			Task{start_time: 4, end_time: 6, cookies: 17},
			Task{start_time: 12, end_time: 13, cookies: 14},
			Task{start_time: 12, end_time: 15, cookies: 23}
		];
		assert_eq!(time_manager(&mut tasks), 84);
	}
}
`

func ex07Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest07, exercise.TurnInFiles[0]); err != nil {
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex07() Exercise.Exercise {
	return Exercise.NewExercise("07", "ex07", []string{"src/lib.rs", "Cargo.toml"}, 20, ex07Test)
}
