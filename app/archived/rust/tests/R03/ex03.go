//go:build ignore
package R03

import (
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var Ex03TestMod = `
#[cfg(test)]
mod shortinette_rust_test_module03_ex03_0001 {
    use std::cell::Cell;

    use super::*;

    #[test]
    fn new_type() {
        #[derive(PartialEq, Eq, Debug)]
        struct NewType;

        impl FortyTwo for NewType {
            fn forty_two() -> Self {
                Self
            }
        }

        assert_eq!(<NewType as FortyTwo>::forty_two(), NewType, "FortyTwo::forty_two() did not return correct result for a custom type");
    }

    #[test]
    fn obvious() {
        assert_eq!(u32::forty_two(), 42);

        String::forty_two();
    }

    // Probably the most hacky test ever.
    // It does not check whether the println! output is correct.
    // But come on, if FortyTwo::forty_two() was called AND the debug formatting
    // was called then it makes no sense to not call println! lol.
    #[test]
    fn forty_two() {
        thread_local! {
            static FT_CALLED: Cell<bool> = const { Cell::new(false) };
            static DEBUG_CALLED: Cell<bool> = const { Cell::new(false) };
        }

        #[derive(PartialEq, Eq)]
        struct NewType;

        impl FortyTwo for NewType {
            fn forty_two() -> Self {
                FT_CALLED.set(true);

                Self
            }
        }

        impl Debug for NewType {
            fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                DEBUG_CALLED.set(true);

                f.write_str("42")
            }
        }

        print_forty_two::<NewType>();

        assert!(FT_CALLED.get());
        assert!(DEBUG_CALLED.get());
    }
}
`

// TODO: this needs to test with different mains
func ex03Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	if err := testutils.AppendStringToFile(Ex03TestMod, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module03_ex03_0001"}, testutils.WithTimeout(5*time.Second))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "ex03", []string{"src/main.rs", "Cargo.toml"}, 10, ex03Test)
}
