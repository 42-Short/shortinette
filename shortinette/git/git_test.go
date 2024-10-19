package git

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
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
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("could not load .env")
	}
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
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("could not load .env")
	}

	if err := AddCollaborator("repo", "ireallydonotthinkthatthisgithubuserexists", "read"); err == nil {
		t.Fatalf("non-existing user should throw an error")
	}
}

func TestAddCollaboratorNonExistingPermission(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("could not load .env")
	}

	if err := AddCollaborator("repo", "winstonallo", "fornicate"); err == nil {
		t.Fatalf("non-existing permission level should throw an error")
	}
}

func TestAddNonExistingFilesToRepo(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("could not load .env: %v", err)
	}

	if err := NewRepo("test", true, "this will be deleted soon"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}

	defer func() {
		if err := deleteRepo("test"); err != nil {
			t.Fatalf("could not delete test repo: %v", err)
		}
	}()

	if err := UploadFiles("test", "don't mind me just breaking code", "foo", "bar"); err == nil {
		t.Fatalf("trying to upload non-existing files to a repo should throw an error")
	}
}
