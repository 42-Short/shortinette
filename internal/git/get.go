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

	repo, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: os.Getenv("GITHUB_USER"),
			Password: os.Getenv("GITHUB_TOKEN"),
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
		return nil, fmt.Errorf("could not open repository: %w", err)
	}
	return repo, nil
}

func cloneOrOpen(repoURL, targetDir string) (*git.Repository, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return cloneRepository(repoURL, targetDir)
	}
	return openRepository(targetDir)
}

func pullLatestChanges(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.Main,
		Auth: &http.BasicAuth{
			Username: os.Getenv("GITHUB_USER"),
			Password: os.Getenv("GITHUB_TOKEN"),
		},
		Progress: os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("could not pull repository: %w", err)
	}

	fmt.Println("repository pulled successfully")
	return nil
}

func get(repoURL, targetDir string) error {
	repo, err := cloneOrOpen(repoURL, targetDir)
	if err != nil {
		return err
	}
	return pullLatestChanges(repo)
}
