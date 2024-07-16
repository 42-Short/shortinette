package git

import (
	"fmt"

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
		logger.Error.Println(err)
		return fmt.Errorf("could not upload %s to repo %s: %w", localFilePath, repoId, err)
	}
	logger.Info.Printf("%s successfully uploaded to %s/%s", localFilePath, repoId, targetFilePath)
	return nil
}
