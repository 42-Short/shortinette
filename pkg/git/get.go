package git

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/42-Short/shortinette/internal/datastructures"
)

func cloneRepository(repoURL, targetDir string, env datastructures.Environment) (*git.Repository, error) {

	repo, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: env.User,
			Password: env.Token,
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

func cloneOrOpen(repoURL, targetDir string, env datastructures.Environment) (*git.Repository, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return cloneRepository(repoURL, targetDir, env)
	}
	return openRepository(targetDir)
}

func pullLatestChanges(repo *git.Repository, targetDir string, env datastructures.Environment) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree for repository in directory %s: %w", targetDir, err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.Main,
		Auth: &http.BasicAuth{
			Username: env.User,
			Password: env.Token,
		},
		Progress: os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error pulling repository %s: %w", targetDir, err)
	}

	fmt.Println("Repository pulled successfully.")
	return nil
}

func get(repoURL, targetDir string, env datastructures.Environment) error {
	repo, err := cloneOrOpen(repoURL, targetDir, env)
	if err != nil {
		return err
	}
	return pullLatestChanges(repo, targetDir, env)
}
