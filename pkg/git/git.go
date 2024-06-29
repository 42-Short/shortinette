package git

import (
	"fmt"
	"os"
)

// Checks if needed enviorment variables are present in the .env file
func CheckRequiredEnvironmentVariables() error {
	vars := map[string]string{
		"GITHUB_USER":         os.Getenv("GITHUB_USER"),
		"GITHUB_TOKEN":        os.Getenv("GITHUB_TOKEN"),
		"GITHUB_ORGANISATION": os.Getenv("GITHUB_ORGANISATION"),
	}
	for key, value := range vars {
		if value == "" {
			return fmt.Errorf("%s environment variable not set", key)
		}
	}
	return nil
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
