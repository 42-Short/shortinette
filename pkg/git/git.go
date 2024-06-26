package git

import (
	"fmt"
	"os"
)

func getCredentials() (string, string, error) {
	username := os.Getenv("GITHUB_USER")
	password := os.Getenv("GITHUB_TOKEN")

	if username == "" || password == "" {
		return username, password, fmt.Errorf("error: GITHUB_USER and/or GITHUB_TOKEN environment variables not set")
	}

	return username, password, nil
}

func Get(repoURL string, targetDirectory string) error {
	if err := get(repoURL, targetDirectory); err != nil {
		return fmt.Errorf("could not get repo: %w", err)
	}
	return nil
}

func Create(name string) error {
	if err := create(name); err != nil {
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}

func AddCollaborator(repo string, name string, permission string) error {
	if err := addCollaborator(repo, name, "push"); err != nil {
		return fmt.Errorf("could not add %s to repo %s: %w", name, repo, err)
	}
	return nil
}