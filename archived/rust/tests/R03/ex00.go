package R03

import (
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var Ex00TestMod = `
#[cfg(test)]
mod shortinette_rust_test_module03_ex00_0001 {
    use super::*;

    #[test]
    #[should_panic]
    fn empty() {
        let empty: Vec<i32> = Vec::new();

        let value = choose(&empty);
        println!(
            "What in the world? This should not be printed! Value: {:?}",
            value
        );
    }

    #[test]
    fn single() {
        let alone = &[0];
        let value = choose(alone);

        assert_eq!(value, &0, "How is choose(&[0]) not returning 0?");
    }

    #[test]
    fn generic() {
        let numbers = [1_u8, 2, 3];
        let _: &u8 = choose(&numbers);

        let slices = ["a", "b", "c"];
        let _: &&str = choose(&slices);

        let bools = [true, false];
        let _: &bool = choose(&bools);

        struct Foo;
        let foos = [Foo, Foo];
        let _: &Foo = choose(&foos);
    }

    #[test]
    fn randomness() {
        let huge: Vec<_> = (0..100_000).collect();

        // Well this is one of the cases where a second grademe could actually
        // pass the exercise.
        // Although it is very unlikely that this will ever happen.
        let value = choose(&huge);
        let value2 = choose(&huge);

        assert_ne!(
            value, value2,
            "choose(&huge) returned the same value twice. Do you really return a random element?"
        );
    }
}
`

func ex00Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	if err := testutils.AppendStringToFile(Ex00TestMod, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module03_ex00_0001"}, testutils.WithTimeout(5*time.Second))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "ex00", []string{"src/main.rs", "Cargo.toml"}, 10, ex00Test)
}
