package git

import (
	"os"
	"testing"
)

func TestMissingRequiredVariables(t *testing.T) {
	os.Unsetenv("GITHUB_ADMIN")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_ORGANISATION")
	os.Unsetenv("GITHUB_EMAIL")

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNonExistingOrga(t *testing.T) {
	os.Setenv("GITHUB_ORGANISATION", "thisorgadoesnotexist")

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}
