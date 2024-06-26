package git

import (
	"fmt"
	"os"
)

func getToken() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	return token, nil
}

// Clone or open the repo & pull the latest changes into targetDirectory
func Get(repoURL string, targetDirectory string) error {
	if err := get(repoURL, targetDirectory); err != nil {
		return err
	}
	return nil
}

// Check if repo exists, if not create it.
func Create(name string) error {
	if err := create(name); err != nil {
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}

// Add a collaborator with the specified permissions to the repo
func AddCollaborator(repo string, name string, permission string) error {
	if err := addCollaborator(repo, name, "push"); err != nil {
		return fmt.Errorf("could not add %s to repo %s: %w", name, repo, err)
	}
	return nil
}
