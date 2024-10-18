package git

import (
	"context"
	"fmt"
	"net/http"
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

// Checks whether `err` is related to the repo already existing.
func isRepoAlreadyExists(err error) (exists bool) {
	if githubErr, ok := err.(*github.ErrorResponse); ok {
		for _, e := range githubErr.Errors {
			if e.Message == "name already exists on this account" {
				return true
			}
		}
	}
	return false
}

// Creates a new repository `name` under the GitHub organisation specified by the
// GITHUB_ORGANISATION environment variable. If `private` is true, the repository's
// visibility will be private.
func NewRepo(name string, private bool, description string) (err error) {
	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	repo := &github.Repository{Name: &name, Private: &private, Description: &description}

	createdRepo, response, err := client.Repositories.Create(ctx, orga, repo)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusUnprocessableEntity {
			if isRepoAlreadyExists(err) {
				fmt.Printf("repo %s already exists under orga %s, skipping\n", name, os.Getenv("GITHUB_ORGANISATION"))
				return nil
			}
		}
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	fmt.Printf("repo created: %s at URL: %s\n", *createdRepo.Name, *createdRepo.HTMLURL)
	return nil
}

// Adds collaborator `collaboratorName` to repo `repoName` (under the GitHub organisation specified by
// the GITHUB_ORGANISATION environment variable) with access level `permission“.
func AddCollaborator(repoName string, collaboratorName string, permission string) (err error) {
	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not add collaborator %s to repo %s: %v", collaboratorName, repoName, err)
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)

	options := &github.RepositoryAddCollaboratorOptions{
		Permission: permission,
	}

	if _, _, err = client.Repositories.AddCollaborator(ctx, orga, repoName, collaboratorName, options); err != nil {
		return fmt.Errorf("could not add collaborator %s to repo %s: %v", collaboratorName, repoName, err)
	}

	fmt.Printf("user %s added to repo %s with %s access\n", collaboratorName, repoName, permission)
	return nil
}
