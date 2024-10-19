package git

import (
	"os"
	"testing"
)

func TestNewRepoMissingRequiredVariables(t *testing.T) {
	os.Unsetenv("ADMIN_GITHUB")
	os.Unsetenv("TOKEN_GITHUB")
	os.Unsetenv("ORGA_GITHUB")
	os.Unsetenv("EMAIL_GITHUB")

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNewRepoNonExistingOrga(t *testing.T) {
	os.Setenv("ORGA_GITHUB", "thisorgadoesnoteist")

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestAddCollaboratorMissingToken(t *testing.T) {
	os.Unsetenv("ADMIN_GITHUB")
	os.Unsetenv("TOKEN_GITHUB")
	os.Unsetenv("ORGA_GITHUB")
	os.Unsetenv("EMAIL_GITHUB")

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

func TestUploadFilesNonExistingFiles(t *testing.T) {
	if err := NewRepo("test", true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}

	defer func() {
		if err := os.RemoveAll("test"); err != nil {
			t.Fatalf("could not delete test repo (local): %v", err)
		}
		if err := deleteRepo("test"); err != nil {
			t.Fatalf("could not delete test repo: %v", err)
		}
	}()

	if err := UploadFiles("test", "don't mind me just breaking code", "foo", "bar"); err == nil {
		t.Fatalf("trying to upload non-existing files to a repo should throw an error")
	}
}

func TestUploadFilesNormalFunctionality(t *testing.T) {
	if err := NewRepo("test", true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}

	defer func() {
		if err := os.RemoveAll("test"); err != nil {
			t.Fatalf("could not delete test repo (local): %v", err)
		}
		if err := deleteRepo("test"); err != nil {
			t.Fatalf("could not delete test repo (remote): %v", err)
		}
	}()

	if err := UploadFiles("test", "don't mind me just breaking code", "git.go"); err != nil {
		t.Fatalf("uploading an existing file should work, something went wrong: %v", err)
	}
}
