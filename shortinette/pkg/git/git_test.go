package git

import (
	"os"
	"testing"
)

func TestNewRepoMissingRequiredVariables(t *testing.T) {
	os.Unsetenv("GITHUB_ADMIN")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_ORGANISATION")
	os.Unsetenv("GITHUB_EMAIL")

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNewRepoNonExistingOrga(t *testing.T) {
	os.Setenv("GITHUB_ORGANISATION", "thisorgadoesnotexist")

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestAddCollaboratorMissingToken(t *testing.T) {
	os.Unsetenv("GITHUB_ADMIN")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_ORGANISATION")
	os.Unsetenv("GITHUB_EMAIL")

	if err := AddCollaborator("repo", "winstonallo", "read"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestAddCollaboratorNonExistingUser(t *testing.T) {
	if err := AddCollaborator("repo", "ireallydonotthinkthatthisgithubuserexists", "read"); err == nil {
		t.Fatalf("non-existing user should throw an error")
	}
}

func TestAddCollaboratorNonExistingPermission(t *testing.T) {
	if err := AddCollaborator("repo", "winstonallo", "fornicate"); err == nil {
		t.Fatalf("non-existing permission level should throw an error")
	}
}