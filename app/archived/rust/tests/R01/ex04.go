//go:build ignore
package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest04 = `
#[cfg(test)]
mod shortinette_tests_rust_0104 {
	use super::*;

	#[test]
	fn test_empty() {
		let mut boxes = [];
		sort_boxes(&mut boxes);
		assert_eq!(boxes.len(), 0, "Failed for an empty list as input");
	}

	#[test]
	fn test_datatype() {
		let mut boxes = [[1u32, 1u32], [2u32, 2u32]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[2u32, 2u32], [1u32, 1u32]], "Failed ordering [[1, 1], [2, 2]]");
	}

	#[test]
	fn test_0() {
		let mut boxes = [[2, 2], [1, 1], [3, 3]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[3, 3], [2, 2], [1, 1]], "Failed ordering [[2, 2], [1, 1], [3, 3]]");
	}

	#[test]
	fn test_1() {
		let mut boxes = [[0, 0], [1, 1], [1, 1], [0, 0]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[1, 1], [1, 1], [0, 0], [0, 0]], "Failed ordering [[0, 0], [1, 1], [1, 1], [0, 0]]");
	}

	#[test]
	fn test_2() {
		let mut boxes = [[0, 1], [1, 1]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[1, 1], [0, 1]], "Failed ordering [[0, 1], [1, 1]]");
	}

	#[test]
	fn test_3() {
		let mut boxes = [[1, 0], [1, 1]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[1, 1], [1, 0]], "Failed ordering [[1, 0], [1, 1]]");
	}

	#[test]
	fn test_4() {
		let mut boxes = [[1, 1], [2, 1]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[2, 1], [1, 1]], "Failed ordering [[1, 1], [2, 1]]");
	}

	#[test]
	fn test_5() {
		let mut boxes = [[1, 1], [1, 2]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[1, 2], [1, 1]], "Failed ordering [[1, 1], [1, 2]]");
	}

	#[test]
	fn test_6() {
		let mut boxes = [[5, 3], [5, 2], [8, 5], [2, 2], [1, 1], [2, 1]];
		sort_boxes(&mut boxes);
		assert_eq!(boxes, [[8, 5], [5, 3], [5, 2], [2, 2], [2, 1], [1, 1]], "Failed ordering [[5, 3], [5, 2], [8, 5], [2, 2], [1, 1], [2, 1]]");
	}

	#[test]
	#[should_panic]
	fn test_7() {
		let mut boxes = [[5, 3], [5, 2], [8, 5], [2, 2], [1, 2], [2, 1]];
		sort_boxes(&mut boxes);
		eprintln!("Did not panic for [[5, 3], [5, 2], [8, 5], [2, 2], [1, 2], [2, 1]]");
	}

	#[test]
	#[should_panic]
	fn test_8() {
		let mut boxes = [[2, 1], [5, 3], [5, 2], [1, 2], [8, 5], [2, 2]];
		sort_boxes(&mut boxes);
		eprintln!("Did not panic for [[2, 1], [5, 3], [5, 2], [1, 2], [8, 5], [2, 2]]");
	}
}
`

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest04, exercise.TurnInFiles[0]); err != nil {
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "ex04", []string{"src/lib.rs", "Cargo.toml"}, 10, ex04Test)
}
