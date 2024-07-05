package R00

import Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"

func ex04Test(exercise *Exercise.Exercise) bool {
	return true
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("EX04", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
