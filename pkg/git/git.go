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
	if err := doGet(repoURL, targetDirectory); err != nil {
		return fmt.Errorf("could not get repo: %w", err)
	}
	return nil
}

func Create(name string) error {
	if err := doCreate(name); err != nil {
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}
