package git

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/v66/github"
)

func TestNewRepoNonExistingOrga(t *testing.T) {
	_, orga, err := requireEnv()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if err := os.Setenv("ORGA_GITHUB", "thisorgadoesnoteist"); err != nil {
		t.Fatalf("error: %v", err)
	}

	defer func() {
		if err := os.Setenv("ORGA_GITHUB", orga); err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	if err := NewRepo("test", true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNewRepoStandardFunctionality(t *testing.T) {
	token, orga, err := requireEnv()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	expectedRepoName := "repository"
	expectedPrivate := true
	expectedDescription := "description"

	if err := NewRepo(expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}

	client := github.NewClient(nil).WithAuthToken(token)
	repo, _, err := client.Repositories.Get(context.Background(), orga, expectedRepoName)
	if err != nil {
		t.Fatalf("could not verify repo creation: %v", err)
	}

	if repo.GetName() != expectedRepoName {
		t.Fatalf("repo name was not set correctly: expected: '%s', got: '%s'", expectedRepoName, repo.GetName())
	}
	if repo.GetPrivate() != expectedPrivate {
		t.Fatalf("repo visibility was not set correctly: expected: private == %t, got: private == %t", expectedPrivate, repo.GetPrivate())
	}
	if repo.GetDescription() != expectedDescription {
		t.Fatalf("repo description was not set correctly: expected: '%s', got: '%s'", expectedDescription, repo.GetDescription())
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
		t.Fatalf("error: %v", err)
	}

	defer func() {
		if err := os.RemoveAll("test"); err != nil {
			t.Fatalf("error: %v", err)
		}
		if err := deleteRepo("test"); err != nil {
			t.Fatalf("error: %v", err)
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

func TestNewReleaseNormalFunctionality(t *testing.T) {
	token, orga, err := requireEnv()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	expectedRepoName := "test"
	expectedTagName := "tag"
	expectedReleaseName := "release"
	expectedBody := "body"

	if err := NewRepo(expectedRepoName, true, "this will be deleted soon"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(expectedRepoName); err != nil {
			t.Fatalf("could not delete test repo (local): %v", err)
		}
		if err := deleteRepo(expectedRepoName); err != nil {
			t.Fatalf("could not delete test repo (remote): %v", err)
		}
	}()

	if err := UploadFiles(expectedRepoName, "initial commit", "git_test.go"); err != nil {
		t.Fatalf("UploadFiles returned an error on initial commit: %v", err)
	}

	if err := NewRelease(expectedRepoName, expectedTagName, expectedReleaseName, expectedBody); err != nil {
		t.Fatalf("NewRelease returned an error on a standard use case: %v", err)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	rel, _, err := client.Repositories.GetLatestRelease(context.Background(), orga, expectedRepoName)
	if err != nil {
		t.Fatalf("could not verify release update: %v", err)
	}

	if rel.GetTagName() != expectedTagName {
		t.Fatalf("release tag name was not set correctly: expected: '%s', got: '%s'", expectedTagName, *rel.TagName)
	}
	if rel.GetName() != expectedReleaseName {
		t.Fatalf("release name was not set correctly: expected: '%s', got '%s'", expectedReleaseName, rel.GetName())
	}
	if rel.GetBody() != expectedBody {
		t.Fatalf("release body was not set correctly: expected: '%s', got: '%s'", expectedBody, rel.GetBody())
	}
}

func TestNewReleaseAlreadyExisting(t *testing.T) {
	expectedRepoName := "test"

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

	if err := UploadFiles(expectedRepoName, "initial commit", "git_test.go"); err != nil {
		t.Fatalf("UploadFiles returned an error on initial commit: %v", err)
	}

	if err := NewRelease("test", "tag", "release", "body"); err != nil {
		t.Fatalf("NewRelease returned an error on a standard use case: %v", err)
	}
}
