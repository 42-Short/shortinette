package git

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v66/github"
)

// Checks for environment variables required to interact with the GitHub API. Returns their values
// if they exist, sets the error's value if not.
func requireEnv() (githubToken string, githubOrga string, err error) {
	missingVars := []string{}

	githubToken = os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		missingVars = append(missingVars, "GITHUB_TOKEN")
	}

	githubOrga = os.Getenv("GITHUB_ORGANISATION")
	if githubOrga == "" {
		missingVars = append(missingVars, "GITHUB_ORGANISATION")
	}

	if len(missingVars) != 0 {
		err = fmt.Errorf("missing environment variable(s): %s", strings.Join(missingVars, ", "))
	}

	return githubToken, githubOrga, err
}

// Creates a new repository under the GitHub organisation specified by the
// GITHUB_ORGANISATION environment variable.
func NewRepo(name string, private bool, description string) (err error) {
	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	repo := &github.Repository{Name: &name, Private: &private, Description: &description}

	createdRepo, _, err := client.Repositories.Create(ctx, orga, repo)
	if err != nil {
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	fmt.Printf("repository created: %s at URL: %s\n", *createdRepo.Name, *createdRepo.HTMLURL)
	return nil
}
