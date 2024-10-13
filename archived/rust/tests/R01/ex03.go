//go:build ignore
package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest03 = `
#[cfg(test)]
mod shortinette_tests_rust_0103 {
	use super::*;

	#[test]
	fn test_lifetimes() {
	    let haystack = [1, 2, 3, 2, 1];
	    let result;

	    {
	        let needle = [2, 3];
	        result = largest_group(&haystack, &needle);
	    }

	    assert_eq!(result, &[2, 3, 2]);
	}

	#[test]
	fn test_datatypes() {
		assert_eq!(largest_group(&[1u32, 2u32, 3u32], &[2u32]), &[2u32]);
	}

	#[test]
	fn test_empty() {
		assert_eq!(largest_group(&[], &[]), &[]);
		assert_eq!(largest_group(&[4, 2], &[]), &[]);
		assert_eq!(largest_group(&[], &[4, 2]), &[]);
	}

	#[test]
	fn test_no_match() {
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[1, 2, 3, 5]), &[]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[5, 1, 2, 3]), &[]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[3, 4, 5, 1]), &[]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[1, 3, 4, 5]), &[]);
		assert_eq!(largest_group(&[1, 2, 3, 4, 5], &[0]), &[]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[15]), &[]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, u32::MAX], &[u32::MAX, 1]), &[]);
	}

	#[test]
	fn test_one_match() {
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[0]), &[0]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[3]), &[3]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[5]), &[5]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5], &[0, 1, 2, 3]), &[0, 1, 2, 3]);
	}

	#[test]
	fn test_duplicates() {
		assert_eq!(largest_group(&[0, 0, 0, 1, 2, 3, 4, 5], &[0]), &[0, 0, 0]);
		assert_eq!(largest_group(&[0, 0, 0, 1, 2, 3, 4, 5], &[0, 0, 0]), &[0, 0, 0]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 3, 3, 4, 5], &[3]), &[3, 3, 3]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 3, 3, 4, 5], &[3, 3, 3]), &[3, 3, 3]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5, 5, 5], &[5]), &[5, 5, 5]);
		assert_eq!(largest_group(&[0, 1, 2, 3, 4, 5, 5, 5], &[5, 5, 5]), &[5, 5, 5]);
	}

	#[test]
	fn test_multiple_matches_different_length() {
		assert_eq!(largest_group(&[9, 5, 3, 5, 9, 5, 3], &[9, 5]), &[5, 9, 5]);
		assert_eq!(largest_group(&[9, 5, 3, 5, 9, 5, 3], &[3, 5]), &[5, 3, 5]);
	}

	#[test]
	fn test_multiple_matches_same_length() {
		assert_eq!(largest_group(&[9, 5, 3, 9, 3, 5], &[3, 5]), &[5, 3]);
		assert_eq!(largest_group(&[9, 5, 3, 0, 9, 3, 5], &[5, 3, 9]), &[9, 5, 3]);
	}
}
`

func ex03Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest03, exercise.TurnInFiles[0]); err != nil {
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "ex03", []string{"src/lib.rs", "Cargo.toml"}, 10, ex03Test)
}
