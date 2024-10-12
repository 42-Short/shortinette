package config

import (
	"testing"
	"time"
)

func TestNewExerciseEmptyExecutablePath(t *testing.T) {
	if _, err := NewExercise("", 10, []string{"foo"}, "bar"); err == nil {
		t.Fatalf("it should not be possible to initialize an exercise with an empty executable path")
	}
}

func TestNewExerciseEmptyTurnInDirectory(t *testing.T) {
	if _, err := NewExercise("foo", 10, []string{"bar"}, ""); err == nil {
		t.Fatalf("it should not be possible to initialize an exercise with an turn in directory")
	}
}

func TestNewExerciseNoAllowedFiles(t *testing.T) {
	if _, err := NewExercise("foo", 10, []string{}, "bar"); err == nil {
		t.Fatalf("it should not be possible to initialize an exercise with no allowed files")
	}
}

func TestNewExerciseNilAllowedFiles(t *testing.T) {
	if _, err := NewExercise("foo", 10, nil, "bar"); err == nil {
		t.Fatalf("allowedFiles cannot be nil")
	}
}

func TestNewExerciseNegativeScore(t *testing.T) {
	if _, err := NewExercise("foo", -10, []string{"bar"}, "bar"); err == nil {
		t.Fatalf("it should not be possible to initialize an exercise with negative score")
	}
}

func TestNewModuleNotEnoughTotalPoints(t *testing.T) {
	ex, _ := NewExercise("foo", 10, []string{"bar"}, "bar")
	exercises := []Exercise{
		0: *ex,
	}
	if _, err := NewModule(exercises, 20); err == nil {
		t.Fatalf("the score of all exercises should add up to the minimum score required to pass (or more)")
	}
}

func TestNewModuleNegativeMinimumScore(t *testing.T) {
	ex, _ := NewExercise("foo", 10, []string{"bar"}, "bar")
	exercises := []Exercise{
		0: *ex,
	}
	if _, err := NewModule(exercises, -20); err == nil {
		t.Fatalf("minimumScore cannot be negative")
	}
}

func TestNewModuleNoExercises(t *testing.T) {
	if _, err := NewModule([]Exercise{}, 0); err == nil {
		t.Fatalf("it should not be possible to initialize a module with no exercises")
	}
}

func TestNewModuleNilExercises(t *testing.T) {
	if _, err := NewModule(nil, 0); err == nil {
		t.Fatalf("exercises cannot be nil")
	}
}

func TestNewShortNegativeDuration(t *testing.T) {
	ex, _ := NewExercise("foo", 10, []string{"bar"}, "bar")
	mod, _ := NewModule([]Exercise{0: *ex}, 10)

	if _, err := NewShort([]Module{0: *mod}, time.Hour*-24, time.Now()); err == nil {
		t.Fatalf("moduleDuration should not be negative")
	}
}

func TestNewShortNoModules(t *testing.T) {
	if _, err := NewShort([]Module{}, time.Hour*24, time.Now()); err == nil {
		t.Fatalf("should not be possible to initialize a Short with no modules")
	}
}

func TestNewShortNilModules(t *testing.T) {
	if _, err := NewShort(nil, time.Hour*24, time.Now()); err == nil {
		t.Fatalf("modules cannot be nil")
	}
}
