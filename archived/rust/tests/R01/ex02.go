package R01

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest02 = `
#[cfg(test)]
mod shortinette_tests_rust_0102 {
	use super::*;

	const fn color_name_validation(color: &[u8; 3]) -> &'static str {
	    match color {
	        [0, 0, 0] => "pure black",
	        [255, 255, 255] => "pure white",
	        [255, 0, 0] => "pure red",
	        [0, 255, 0] => "pure green",
	        [0, 0, 255] => "pure blue",
	        [128, 128, 128] => "perfect grey",
	        [0..=30, 0..=30, 0..=30] => "almost black",
	        [129..=255, 0..=127, 0..=127] => "redish",
	        [0..=127, 129..=255, 0..=127] => "greenish",
	        [0..=127, 0..=127, 129..=255] => "blueish",
	        _ => "unknown"
	    }
	}

	fn test_color(color: &[u8; 3]) {
		let name_of_the_best_color;
		let expected_color;

		{
			name_of_the_best_color = color_name(color);
			expected_color = color_name_validation(color);
		}

		assert_eq!(name_of_the_best_color, expected_color);
	}

	#[test]
	fn test_0() {
		let test_values = [0, 1, 29, 30, 31, 126, 127, 128, 129, 254, 255];
		for r in test_values {
			for g in test_values {
				for b in test_values {
					test_color(&[r, g, b]);
				}
			}
		}
	}
}
`

func constKeywordCheck(exercise *Exercise.Exercise) Exercise.Result {
	content, err := os.ReadFile(exercise.TurnInFiles[0])
	if err != nil {
		return Exercise.InternalError(fmt.Sprintf("error reading file: %v", err.Error()))
	}
	pattern := `const.*\bfn\scolor_name\b`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	if matches := regex.FindAll(content, -1); len(matches) == 0 {
		return Exercise.CompilationError("color_name function must be declared as const")
	}
	return Exercise.Passed("")
}

func ex02Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0, "if": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if result := constKeywordCheck(exercise); !result.Passed {
		return result
	}
	if err := testutils.AppendStringToFile(cargoTest02, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("02", "ex02", []string{"src/lib.rs", "Cargo.toml"}, 10, ex02Test)
}
