package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func cloneOrOpen(repoURL string, targetDir string) (*git.Repository, error) {
	var repo *git.Repository

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		repo, err = git.PlainClone(targetDir, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, fmt.Errorf("error cloning repository %s to directory %s: %w", repoURL, targetDir, err)
		}
	} else {
		repo, err = git.PlainOpen(targetDir)
		if err != nil {
			return nil, fmt.Errorf("error opening repository in directory %s: %w", targetDir, err)
		}
	}
	return repo, nil
}

func doGet(repoURL string, targetDir string) error {

	var repo *git.Repository
	var err error

	repo, err = cloneOrOpen(repoURL, targetDir)
	if err != nil {
		return err
	}

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
		Auth: &githttp.BasicAuth{
			Username: username,
			Password: password,
		},
		Progress: os.Stdout,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error pulling repository %s: %w", repoURL, err)
	}

	fmt.Println("Repository pulled successfully.")
	return nil
}
