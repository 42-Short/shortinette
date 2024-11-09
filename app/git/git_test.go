package git

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v66/github"
	"github.com/google/uuid"
)

func cleanup(t *testing.T, repoName string) {
	if err := os.RemoveAll(repoName); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	if err := deleteRepo(repoName); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}

func TestNewRepoNonExistingOrga(t *testing.T) {
	repoName := uuid.New().String()

	_, orga, _, err := requireEnv()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if err := os.Setenv("ORGA_GITHUB", "thisorgadoesnoteist"); err != nil {
		t.Fatalf("error: %v", err)
	}

	defer func() {
		cleanup(t, repoName)
		if err := os.Setenv("ORGA_GITHUB", orga); err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	if err := NewRepo(repoName, true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNewRepoStandardFunctionality(t *testing.T) {
	t.Parallel()

	token, orga, _, err := requireEnv()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	expectedRepoName := uuid.New().String()
	expectedPrivate := true
	expectedDescription := "description"

	if err := NewRepo(expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, expectedRepoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

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

func TestNewRepoAlreadyExisting(t *testing.T) {
	t.Parallel()

	expectedRepoName := uuid.New().String()
	expectedPrivate := true
	expectedDescription := "description"

	if err := NewRepo(expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, expectedRepoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := NewRepo(expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo should not error on already exsiting repos: %v", err)
	}
}

func TestAddCollaboratorNonExistingUser(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "idc"); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, repoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := AddCollaborator("repo", "ireallydonotthinkthatthisgithubuserexists", "read"); err == nil {
		t.Fatalf("non-existing user should throw an error")
	}
}

func TestAddCollaboratorNonExistingPermission(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "idc"); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, repoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := AddCollaborator("repo", "winstonallo", "fornicate"); err == nil {
		t.Fatalf("non-existing permission level should throw an error")
	}
}

func TestUploadFilesNonExistingFiles(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("error: %v", err)
	}
	defer cleanup(t, repoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := UploadFiles(repoName, "don't mind me just breaking code", "main", false, "foo", "bar"); err == nil {
		t.Fatalf("trying to upload non-existing files to a repo should throw an error")
	}
}

func TestUploadFilesNormalFunctionality(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, repoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := UploadFiles(repoName, "don't mind me just breaking code", "main", false, "git.go", "git_test.go"); err != nil {
		t.Fatalf("uploading an existing file should work, something went wrong: %v", err)
	}

	if err := Clone(repoName); err != nil {
		t.Fatalf("could not verify file upload: %v", err)
	}

	content, err := os.ReadDir(repoName)
	if err != nil {
		t.Fatalf("could not verify file upload: %v", err)
	}

	found := 0
	for _, path := range content {
		if path.Name() == "git.go" || path.Name() == "git_test.go" {
			found += 1
		}
	}

	if found != 2 {
		t.Fatalf("expected 'git.go' and 'git_test.go' to be uploaded, did not find all of them")
	}
}

func TestUploadFilesNonExistingBranch(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, repoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := UploadFiles(repoName, "don't mind me just breaking code", "thisbranchdoesnotexist", false, "git.go", "git_test.go"); err == nil {
		t.Fatalf("UploadFiles should return an error when trying to push to unexisting branch")
	}
}

func TestUploadFilesNonDefaultBranch(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, repoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := UploadFiles(repoName, "don't mind me just breaking code", "thisbranchshouldbecreated", true, "git.go", "git_test.go"); err != nil {
		t.Fatalf("UploadFiles should be able to create a new branch when needed")
	}

	if err := Clone(repoName); err != nil {
		t.Fatalf("could not verify file upload: %v", err)
	}

	content, err := os.ReadDir(repoName)
	if err != nil {
		t.Fatalf("could not verify file upload: %v", err)
	}

	found := 0
	for _, path := range content {
		if path.Name() == "git.go" || path.Name() == "git_test.go" {
			found += 1
		}
	}

	if found != 2 {
		t.Fatalf("expected 'git.go' and 'git_test.go' to be uploaded, did not find all of them")
	}
}

func TestNewReleaseNormalFunctionality(t *testing.T) {
	t.Parallel()

	token, orga, _, err := requireEnv()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	expectedRepoName := uuid.New().String()
	expectedTagName := "tag"
	expectedReleaseName := "release"
	expectedBody := "body"

	if err := NewRepo(expectedRepoName, true, "this will be deleted soon"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, expectedRepoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := UploadFiles(expectedRepoName, "initial commit", "main", false, "git_test.go"); err != nil {
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
	t.Parallel()

	expectedRepoName := uuid.New().String()

	if err := NewRepo(expectedRepoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, expectedRepoName)

	time.Sleep(5 * time.Second) // Generating templates takes a few seconds

	if err := UploadFiles(expectedRepoName, "initial commit", "main", false, "git_test.go"); err != nil {
		t.Fatalf("UploadFiles returned an error on initial commit: %v", err)
	}

	if err := NewRelease(expectedRepoName, "tag", "release", "body"); err != nil {
		t.Fatalf("NewRelease returned an error on a standard use case: %v", err)
	}

	if err := NewRelease(expectedRepoName, "tag", "release", "body"); err == nil {
		t.Fatalf("duplicate tag names should return an error")
	}
}

func TestNewReleaseNonExistingrepo(t *testing.T) {
	t.Parallel()

	repoName := uuid.New().String()

	if err := NewRepo(repoName, true, "idc"); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, repoName)

	if err := NewRelease("thisrepodoesnotexist", "tag", "release", "body"); err == nil {
		t.Fatalf("NewRelease did not return any error when trying to add a release to a non-existing repo")
	}
}
