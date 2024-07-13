# Documentation for `git` Package
## Overview
The `git` package provides utilities to interact with the GitHub repos. It supports creating/cloning repos, pulling changes and adding collaborators.

## Package Structure
* **create.go**: Repo creation
* **clone.go**: Cloning
* **update.go**: Collaborators updates, file uploads
* **git.go**: High-level wrappers for all the above

## Usage
### Example
```go
package main

import (
	"log"
	"pkg/git"
)

func main() {
	repoName := "test-repo"
	repoURL := "https://github.com/42-Short/test-repo"
	targetDir := "local-repo"

	// Create a new repository
	if err := git.Create(repoName); err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	} else {
		log.Println("Repository created successfully.")
	}

	// Clone or open the repository and pull the latest changes
	if err := git.Get(repoURL, targetDir); err != nil {
		log.Fatalf("Failed to get repository: %v", err)
	} else {
		log.Println("Repository cloned/pulled successfully.")
	}

	// Add a collaborator to the repository
	if err := git.AddCollaborator(repoName, "collaborator-username", "push"); err != nil {
		log.Fatalf("Failed to add collaborator: %v", err)
	} else {
		log.Println("Collaborator added successfully.")
	}
}
```
## Environment Variables
The following environment variables need to be set:
* **`GITHUB_ADMIN`**: GitHub username for authentication.
* **`GITHUB_TOKEN`**: GitHub personal access token for authentication

## Error Handling
All functions return errors using Go's error handling pattern. Errors are wrapped with additional context using fmt.Errorf.
