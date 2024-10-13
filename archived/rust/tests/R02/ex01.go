package R02

import (
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

var cargoTestModAsString01 = `

#[cfg(test)]
mod shortinette_tests_rust_0201 {
    use super::*;

    #[test]
    fn test_point_new() {
        let point = Point::new(12.0, 42.0);
        assert_eq!(point.x, 12.0);
        assert_eq!(point.y, 42.0);
    }

    #[test]
    fn test_point_zero() {
        let point = Point::zero();
        assert_eq!(point.x, 0.0);
        assert_eq!(point.y, 0.0);
    }

    #[test]
    fn test_point_distance_to_zero() {
        assert_eq!(Point::new(6.9, 42.0).distance(&Point::zero()), 42.563012111);
    }

    #[test]
    fn test_point_distance_same_point() {
        assert_eq!(Point::new(10.0, 10.0).distance(&Point::new(10.0, 10.0)), 0.0);
    }

    #[test]
    fn test_point_distance_positive_negative() {
        assert_eq!(Point::new(-5.0, 5.0).distance(&Point::new(5.0, -5.0)), 14.1421356237);
    }

    #[test]
    fn test_point_distance_large_values() {
        assert_eq!(Point::new(10000.0, 30000.0).distance(&Point::new(20000.0, 40000.0)), 14142.135623730951);
    }

    #[test]
    fn test_point_distance_fractional() {
        assert_eq!(Point::new(0.5, 0.5).distance(&Point::new(1.5, 1.5)), 1.41421356237);
    }
    
    fn setup_point_translate_and_assert(initial_point: Point, translation: Point, expected: Point) {
        let mut point = initial_point;
        point.translate(translation.x, translation.y);
        assert_eq!(point.x, expected.x);
        assert_eq!(point.y, expected.y);
    }
    
    #[test]
    fn test_point_translate_positive() {
        setup_point_translate_and_assert(Point::new(1.0, 1.0), Point::new(5.0, 3.0), Point::new(6.0, 4.0));
    }
    
    #[test]
    fn test_point_translate_negative() {
        setup_point_translate_and_assert(Point::new(5.0, 5.0), Point::new(-2.0, -3.0), Point::new(3.0, 2.0));
    }
    
    #[test]
    fn test_point_translate_zero() {
        setup_point_translate_and_assert(Point::new(2.0, 3.0), Point::new(0.0, 0.0), Point::new(2.0, 3.0));
    }
    
    #[test]
    fn test_point_translate_to_zero() {
        setup_point_translate_and_assert(Point::new(3.0, 4.0), Point::new(-3.0, -4.0), Point::new(0.0, 0.0));
    }
    
    #[test]
    fn test_point_translate_fractional() {
        setup_point_translate_and_assert(Point::new(1.5, 2.5), Point::new(0.5, 0.5), Point::new(2.0, 3.0));
    }
}

`

var clippyTomlAsString01 = ``

func ex01Test(exercise *Exercise.Exercise) Exercise.Result {
	return runDefaultTest(exercise, cargoTestModAsString01, clippyTomlAsString01, map[string]int{"unsafe": 0})
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"src/lib.rs", "Cargo.toml"}, 10, ex01Test)
}
