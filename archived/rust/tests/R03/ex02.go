package R03

import (
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var Ex02TestMod = `
#[cfg(test)]
mod shortinette_rust_test_module03_ex02_0001 {
    use super::*;

    #[test]
    fn display() {
        let john = John;

        assert_eq!(format!("{}", john), "Hey! I'm John.");
    }

    #[test]
    fn width() {
        let john = John;

        assert_eq!(
            format!("|{:<30}|", john),
            "|Hey! I'm John.                |"
        );

        assert_eq!(
            format!("|{:>30}|", john),
            "|                Hey! I'm John.|"
        );

        assert_eq!(
            format!("|{:^30}|", john),
            "|        Hey! I'm John.        |"
        );
    }

    #[test]
    fn precision() {
        let john = John;

        assert_eq!(format!("{john:.6}"), "Hey! I");
        assert_eq!(format!("{john:.100}"), "Hey! I'm John.");

        assert_eq!(format!("{john:.0}"), "Don't try to silence me!");
    }

    #[test]
    fn debug() {
        let john = John;

        assert_eq!(format!("{john:?}"), "John, the man himself.");
    }

    #[test]
    fn debug_alternate() {
        let john = John;

        assert_eq!(
            format!("{john:#?}"),
            "John, the man himself. He's handsome AND formidable."
        );
    }

    #[test]
    fn binary() {
        let john = John;

        assert_eq!(format!("{john:b}"), "Bip Boop?");
    }
}
`

func ex02Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	if err := testutils.AppendStringToFile(Ex02TestMod, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module03_ex02_0001"}, testutils.WithTimeout(5*time.Second))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("02", "ex02", []string{"src/main.rs", "Cargo.toml"}, 10, ex02Test)
}
