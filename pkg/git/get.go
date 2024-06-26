package git

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func cloneRepository(repoURL, targetDir string) (*git.Repository, error) {
	username, token, err := getCredentials()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})
	if err != nil {
		if err.Error() == "remote repository is empty" {
			return nil, errors.NewSubmissionError(errors.ErrEmptyRepo, "gg lol")
		}
		return nil, fmt.Errorf("could not clone repository %s to directory %s: %w", repoURL, targetDir, err)
	}
	return repo, nil
}

func openRepository(targetDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(targetDir)
	if err != nil {
		return nil, fmt.Errorf("error opening repository in directory %s: %w", targetDir, err)
	}
	return repo, nil
}

func cloneOrOpen(repoURL, targetDir string) (*git.Repository, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return cloneRepository(repoURL, targetDir)
	}
	return openRepository(targetDir)
}

func getCredentials() (string, string, error) {
	username := os.Getenv("GITHUB_USER")
	token := os.Getenv("GITHUB_TOKEN")
	if username == "" || token == "" {
		return "", "", fmt.Errorf("error: GITHUB_USER and/or GITHUB_TOKEN environment variables not set")
	}
	return username, token, nil
}

func pullLatestChanges(repo *git.Repository, targetDir string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree for repository in directory %s: %w", targetDir, err)
	}

	username, password, err := getCredentials()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.Main,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
		Progress: os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error pulling repository %s: %w", targetDir, err)
	}

	fmt.Println("Repository pulled successfully.")
	return nil
}

func get(repoURL, targetDir string) error {
	repo, err := cloneOrOpen(repoURL, targetDir)
	if err != nil {
		return err
	}
	return pullLatestChanges(repo, targetDir)
}
