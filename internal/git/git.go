package git

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/logger"
)

// Clone or open the repo & pull the latest changes into targetDirectory
func Get(repoURL string, targetDirectory string) error {
	if err := get(repoURL, targetDirectory); err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

// Check if repo exists, if not create it.
func Create(name string) error {
	if err := create(name); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}

// Add a collaborator with the specified permissions to the repo
func AddCollaborator(repoId string, name string, permission string) error {
	if err := addCollaborator(repoId, name, permission); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not add %s to repo %s: %w", name, repoId, err)
	}
	return nil
}

// Add/Update a file on a repository
//
// @params:
//   - repoId: The name of the organisation repository
//   - localFilePath: The source file whose content is to be uploaded
//   - targetFilePath: The file to be created/updated on the remote
func UploadFile(repoId string, localFilePath string, targetFilePath string) error {
	if err := uploadFile(repoId, localFilePath, targetFilePath); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not upload %s to repo %s: %w", localFilePath, repoId, err)
	}
	return nil
}

func RemoveCollaborator(repoId string, collaborator string) error {
	if err := removeCollaborator(repoId, collaborator); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not remove collaborator %s from repo %s: %w", collaborator, repoId, err)
	}
	return nil
}