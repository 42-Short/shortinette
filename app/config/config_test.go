package config

import (
	"os"
	"testing"
	"time"
)

func TestFetchEnvVariables(t *testing.T) {
	c := NewConfig([]Participant{}, []Module{}, 0, time.Now())
	if err := c.FetchEnvVariables(); err != nil {
		t.Fatalf("failed to fetch env variables")
	}
}

func TestFetchEnvVariablesMissing(t *testing.T) {
	os.Setenv("ORGA_GITHUB", "")
	c := NewConfig([]Participant{}, []Module{}, 0, time.Now())
	if err := c.FetchEnvVariables(); err == nil {
		t.Fatalf("missing env vars should throw an error")
	}
}

func TestNewParticipantsNonExistingJsonPath(t *testing.T) {
	c := NewConfig([]Participant{}, []Module{}, 0, time.Now())
	if err := c.FetchParticipants("foo"); err == nil {
		t.Fatalf("a non-existing participants list should throw an error")
	}
}

func TestNewParticipantsMalformedJson(t *testing.T) {
	c := NewConfig([]Participant{}, []Module{}, 0, time.Now())
	if err := c.FetchParticipants("config/malformed.json"); err == nil {
		t.Fatalf("a malformed participants list should throw an error")
	}
}

func TestNewParticipantsEmptyList(t *testing.T) {
	c := NewConfig([]Participant{}, []Module{}, 0, time.Now())
	if err := c.FetchParticipants("config/empty.json"); err == nil {
		t.Fatalf("an empty participants list should throw an error")
	}
}

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
