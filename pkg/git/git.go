package git

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/datastructures"
)

// Retrieves environment variables related to GitHub configuration
func GetEnvironment() (datastructures.Environment, error) {
	env := datastructures.Environment{}

	if env.User = os.Getenv("GITHUB_USER"); env.User == "" {
		return env, fmt.Errorf("GITHUB_USER environment variable not set")
	} else if env.Token = os.Getenv("GITHUB_TOKEN"); env.Token == "" {
		return env, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	} else if env.Organisation = os.Getenv("GITHUB_ORGANISATION"); env.Organisation == "" {
		return env, fmt.Errorf("GITHUB_ORGANISATION environment variable not set")
	}
	return env, nil
}

// Clone or open the repo & pull the latest changes into targetDirectory
func Get(repoURL string, targetDirectory string, env datastructures.Environment) error {
	if err := get(repoURL, targetDirectory, env); err != nil {
		return err
	}
	return nil
}

// Check if repo exists, if not create it.
func Create(name string, env datastructures.Environment) error {
	if err := create(name, env); err != nil {
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}

// Add a collaborator with the specified permissions to the repo
func AddCollaborator(repo string, name string, permission string, env datastructures.Environment) error {
	if err := addCollaborator(repo, name, "push", env); err != nil {
		return fmt.Errorf("could not add %s to repo %s: %w", name, repo, err)
	}
	return nil
}
