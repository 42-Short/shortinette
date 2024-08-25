package short

import (
	"os"
	"testing"
	"time"

	"github.com/42-Short/shortinette/pkg/db"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
)

func TestPrematureGradingAttempt(t *testing.T) {
	repo := db.Repository{
		ID:              "foo",
		FirstAttempt:    false,
		LastGradingTime: time.Now().Add(-10 * time.Minute),
		WaitingTime:     15 * time.Minute,
		Score:           0,
		Attempts:        0,
	}
	os.Unsetenv("DEV_MODE")
	if err := checkPrematureGradingAttempt(repo); err == nil {
		t.Fatal("premature rading went through")
	}
}

func TestGradeModuleEmptyInput(t *testing.T) {
	if err := GradeModule(Module.Module{}, "foo"); err == nil {
		t.Fatal("GradeModule runs with empty Module struct")
	}
}
