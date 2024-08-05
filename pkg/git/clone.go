package git

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func cloneRepository(repoURL, targetDir string) (err error) {
	_, err = git.PlainClone(targetDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: os.Getenv("GITHUB_ADMIN"),
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
	if err != nil {
		if err.Error() == "remote repository is empty" {
			return errors.NewSubmissionError(errors.ErrEmptyRepo, "gg lol")
		}
		return fmt.Errorf("could not clone repository %s to directory %s: %w", repoURL, targetDir, err)
	}
	return nil
}

func clone(repoURL, targetDir string) (err error) {
	if _, err = os.Stat(targetDir); os.IsNotExist(err) {
		if err = cloneRepository(repoURL, targetDir); err != nil {
			return err
		}
	}
	return nil
}

func get(repoURL, targetDir string) (err error) {
	err = clone(repoURL, targetDir)
	if err != nil {
		return err
	}
	logger.Info.Printf("repository %s cloned", repoURL)
	return nil
}
