package git

import (
	"fmt"
)

// Clone or open the repository and pull the latest changes into targetDirectory
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

// Add a collaborator with the specified permissions to the repository
func AddCollaborator(repo string, name string, permission string) error {
	if err := addCollaborator(repo, name, "push"); err != nil {
		return fmt.Errorf("could not add %s to repo %s: %w", name, repo, err)
	}
	return nil
}