package git

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/42-Short/shortinette/internal/logger"
)

// Clone a GitHub repo into targetDirectory.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func Clone(repoURL string, targetDirectory string) error {
	if err := get(repoURL, targetDirectory); err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

// Check if repo exists, if not create it under the configured organisation.
// Also adds a webhook for easy recording of repository activity.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func Create(name string) error {
	if err := create(name); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}

// Add a collaborator with the specified permissions to the repo.
//
//   - repoId: name of the organisation repository
//   - username: GitHub username of the collaborator
//   - permission: access level to be given to the user
//
// NOTE: Using this function will overwrite the user's previous rights - use test
// accounts, or you might lock yourself out of your repos.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func AddCollaborator(repoId string, username string, permission string) error {
	if err := addCollaborator(repoId, username, permission); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not add %s to repo %s: %w", username, repoId, err)
	}
	return nil
}

// Add/Update a file on a repository
//
//   - repoId: name of the organisation repository
//   - localFilePath: source file to be uploaded
//   - targetFilePath: file to be created/updated on the remote
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func UploadFile(repoId string, localFilePath string, targetFilePath string, commitMessage string, branch string) error {
	if err := uploadFile(repoId, localFilePath, targetFilePath, commitMessage, branch); err != nil {
		return fmt.Errorf("could not upload %s to repo %s: %w", localFilePath, repoId, err)
	}
	logger.Info.Printf("uploaded %s to repo %s", localFilePath, repoId)
	return nil
}

func NewRelease(repoId string, tagName string, releaseName string, graded bool) error {
	if err := newRelease(repoId, tagName, releaseName, graded); err != nil {
		return err
	}
	logger.Info.Printf("added new release '%s' to %s", releaseName, repoId)
	return nil
}

func GetLatestScore(repoid string) int {
	_, releaseName, _, err := getLatestRelease(repoid)
	if err != nil {
		logger.Error.Printf("failed getting the latest release for repo %s: %v", repoid, err)
		return 0
	}

	nameParts := strings.Split(releaseName, "/")
	if len(nameParts) == 0 {
		logger.Error.Printf("invalid release name format: %s", releaseName)
		return 0
	}
	score, err := strconv.Atoi(nameParts[0])
	if err != nil {
		score = 0
	}

	return score
}
