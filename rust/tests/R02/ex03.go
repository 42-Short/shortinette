package R02

import (
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

var cargoTestModAsString03 = `

#[cfg(test)]
mod shortinette_tests_rust_0203 {
    use super::*;

    #[test]
    fn test_clone_trait() {
        let instance = MyType::default();
        assert_eq!(instance, instance.clone());
    }

    #[test]
    fn test_partial_eq_trait() {
        assert_eq!(MyType::default(), MyType::default());
    }

    #[test]
    fn test_partial_ord_trait() {
        let instance1 = MyType::default();
        let instance2 = MyType::default();
        assert!(instance1 <= instance2 && instance1 >= instance2);
    }

    #[test]
    fn test_debug_trait() {
        format!("{:?}", MyType::default());
    }
}

`

var clippyTomlAsString03 = ``

func ex03Test(exercise *Exercise.Exercise) Exercise.Result {
	return runDefaultTest(exercise, cargoTestModAsString03, clippyTomlAsString03, map[string]int{"impl": 0, "unsafe": 0})
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "ex03", []string{"src/main.rs", "Cargo.toml"}, 10, ex03Test)
}
