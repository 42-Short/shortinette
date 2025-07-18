package git

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v66/github"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	orga         string
	token        string
	templateRepo string
	basePath     string
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env")

	orga = os.Getenv("ORGA_GITHUB")
	token = os.Getenv("TOKEN_GITHUB")
	templateRepo = os.Getenv("TEMPLATE_REPO")
	basePath = os.Getenv("BASE_PATH")
	code := m.Run()
	os.Exit(code)
}

func cleanup(t *testing.T, gh *GithubService, repoName string) {
	if err := os.RemoveAll(repoName); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	if err := gh.deleteRepo(repoName); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}

func TestNewRepoNonExistingOrga(t *testing.T) {
	gh := NewGithubService(token, "thisorgadoesnoteist", basePath)
	repoName := uuid.New().String()

	defer cleanup(t, gh, repoName)

	if err := gh.NewRepo(templateRepo, repoName, true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNewRepoMissingTemplateRepoEnvironmentVariable(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo("", repoName, true, "this should not be created"); err == nil {
		t.Fatalf("missing environment variables should throw an error")
	}
}

func TestNewRepoNonExistingTemplate(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	expectedRepoName := uuid.New().String()

	if err := gh.NewRepo("thistemplatedoesnotexist", expectedRepoName, true, "expectedDescription"); err == nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
}

func TestNewRepoStandardFunctionality(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	expectedRepoName := uuid.New().String()
	expectedPrivate := true
	expectedDescription := "description"

	if err := gh.NewRepo(templateRepo, expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, gh, expectedRepoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

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
	gh := NewGithubService(token, orga, basePath)
	expectedRepoName := uuid.New().String()
	expectedPrivate := true
	expectedDescription := "description"

	if err := gh.NewRepo(templateRepo, expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, gh, expectedRepoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.NewRepo(templateRepo, expectedRepoName, expectedPrivate, expectedDescription); err != nil {
		t.Fatalf("NewRepo should not error on already existing repos: %v", err)
	}
}

func TestAddCollaboratorNonExistingUser(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, repoName, true, "idc"); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, gh, repoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.AddCollaborator("repo", "ireallydonotthinkthatthisgithubuserexists", "read"); err == nil {
		t.Fatalf("non-existing user should throw an error")
	}
}

func TestAddCollaboratorNonExistingPermission(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, repoName, true, "idc"); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, gh, repoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.AddCollaborator("repo", "winstonallo", "fornicate"); err == nil {
		t.Fatalf("non-existing permission level should throw an error")
	}
}

func TestUploadFilesNonExistingFiles(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("error: %v", err)
	}
	defer cleanup(t, gh, repoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.UploadFiles(repoName, "don't mind me just breaking code", "main", false, "foo", "bar"); err == nil {
		t.Fatalf("trying to upload non-existing files to a repo should throw an error")
	}
}

func TestUploadFilesNormalFunctionality(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, gh, repoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.UploadFiles(repoName, "don't mind me just breaking code", "main", false, "git.go", "git_test.go"); err != nil {
		t.Fatalf("uploading an existing file should work, something went wrong: %v", err)
	}

	if err := gh.Clone(repoName); err != nil {
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
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, repoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, gh, repoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.UploadFiles(repoName, "don't mind me just breaking code", "thisbranchdoesnotexist", false, "git.go", "git_test.go"); err == nil {
		t.Fatalf("UploadFiles should return an error when trying to push to unexisting branch")
	}
}

func TestNewReleaseNormalFunctionality(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	expectedRepoName := uuid.New().String()
	expectedTagName := "tag"
	expectedReleaseName := "release"
	expectedBody := "body"

	if err := gh.NewRepo(templateRepo, expectedRepoName, true, "this will be deleted soon"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, gh, expectedRepoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.UploadFiles(expectedRepoName, "initial commit", "main", false, "git_test.go"); err != nil {
		t.Fatalf("UploadFiles returned an error on initial commit: %v", err)
	}

	if err := gh.NewRelease(expectedRepoName, expectedTagName, expectedReleaseName, expectedBody); err != nil {
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
	gh := NewGithubService(token, orga, basePath)
	expectedRepoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, expectedRepoName, true, "this will be deleted soon_GITHUB"); err != nil {
		t.Fatalf("could not create test repo: %v", err)
	}
	defer cleanup(t, gh, expectedRepoName)

	time.Sleep(3 * time.Second) // Generating templates takes a few seconds

	if err := gh.UploadFiles(expectedRepoName, "initial commit", "main", false, "git_test.go"); err != nil {
		t.Fatalf("UploadFiles returned an error on initial commit: %v", err)
	}

	if err := gh.NewRelease(expectedRepoName, "tag", "release", "body"); err != nil {
		t.Fatalf("NewRelease returned an error on a standard use case: %v", err)
	}

	if err := gh.NewRelease(expectedRepoName, "tag", "release", "body"); err == nil {
		t.Fatalf("duplicate tag names should return an error")
	}
}

func TestNewReleaseNonExistingrepo(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)
	repoName := uuid.New().String()

	if err := gh.NewRepo(templateRepo, repoName, true, "idc"); err != nil {
		t.Fatalf("NewRepo returned an error on a standard use case: %v", err)
	}
	defer cleanup(t, gh, repoName)

	if err := gh.NewRelease("thisrepodoesnotexist", "tag", "release", "body"); err == nil {
		t.Fatalf("NewRelease did not return any error when trying to add a release to a non-existing repo")
	}
}

func TestDoesAccountExistNonExisting(t *testing.T) {
	found, err := DoesAccountExist("thisuserdoesnotexist_42424242424242424242424000")
	require.NoError(t, err)
	assert.Equal(t, found, false, "DoesAccountExist returned true on an invalid user")
}

func TestDoesAccountExistExisting(t *testing.T) {
	found, err := DoesAccountExist("winstonallo")
	require.NoError(t, err)
	assert.Equal(t, found, true, "DoesAccountExist returned false on a valid user")
}

func TestNewTemplateRepo(t *testing.T) {
	gh := NewGithubService(token, orga, basePath)

	if _, err := gh.CreateModuleTemplate(0); err != nil {
		t.Fatalf("NewTemplateRepo failed on a standard use case: %v", err)
	}

	defer cleanup(t, gh, "module-00-template")
}
