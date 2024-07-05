package R00

import Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"

func ex03Test(exercise *Exercise.Exercise) bool {
	return true
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("EX03", "ex03", "fizzbuzz.rs", "program", "", []string{"println"}, nil, map[string]int{"match": 1, "for": 0}, ex03Test)
}
