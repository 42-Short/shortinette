// `clone.go (git package)` provides functions for cloning GitHub repositories using the go-git library.
package git

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// cloneRepository clones a GitHub repository from the specified repoURL into the targetDir.
//
//   - repoURL: the URL of the repository to clone
//   - targetDir: the directory where the repository should be cloned
//
// Returns an error if the cloning process fails.
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
			return fmt.Errorf("empty repository")
		}
		return fmt.Errorf("could not clone repository %s to directory %s: %w", repoURL, targetDir, err)
	}
	return nil
}

// clone checks if the target directory exists. If it does not, it clones the repository
// from the specified repoURL into the targetDir.
//
//   - repoURL: the URL of the repository to clone
//   - targetDir: the directory where the repository should be cloned
//
// Returns an error if the cloning process fails.
func clone(repoURL, targetDir string) (err error) {
	if _, err = os.Stat(targetDir); os.IsNotExist(err) {
		if err = cloneRepository(repoURL, targetDir); err != nil {
			return err
		}
	}
	return nil
}

// get clones the repository from the specified repoURL into the targetDir. If the
// repository is successfully cloned, a log message is printed.
//
//   - repoURL: the URL of the repository to clone
//   - targetDir: the directory where the repository should be cloned
//
// Returns an error if the cloning process fails.
func get(repoURL, targetDir string) (err error) {
	err = clone(repoURL, targetDir)
	if err != nil {
		return err
	}
	logger.Info.Printf("repository %s cloned", repoURL)
	return nil
}
