package git

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/42-Short/shortinette/pkg/logger"
)

func sendHTTPRequest(request *http.Request) (response *http.Response, err error) {
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		return response, fmt.Errorf("request failed: %s, %s", response.Status, body)
	}
	return response, nil
}

func createHTTPRequest(method, url, token string, body []byte) (*http.Request, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// Clone a GitHub repo from repoURL into targetDirectory.
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
//   - repoID: name of the organisation repository
//   - username: GitHub username of the collaborator
//   - permission: access level to be given to the user
//
// NOTE: Using this function will overwrite the user's previous rights - use test
// accounts, or you might lock yourself out of your repos.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func AddCollaborator(repoID string, username string, permission string) error {
	if err := addCollaborator(repoID, username, permission); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not add %s to repo %s: %w", username, repoID, err)
	}
	return nil
}

// Add/Update a file on a repository
//
//   - repoID: name of the organisation repository
//   - localFilePath: source file to be uploaded
//   - targetFilePath: file to be created/updated on the remote
//   - commitMessage: the message which will be added to the commit
//   - branch: the branch the data is to be pushed to
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func UploadFile(repoID string, localFilePath string, targetFilePath string, commitMessage string, branch string) error {
	if err := uploadFile(repoID, localFilePath, targetFilePath, commitMessage, branch); err != nil {
		return fmt.Errorf("could not upload %s to repo %s: %w", localFilePath, repoID, err)
	}
	logger.Info.Printf("uploaded %s to repo %s", localFilePath, repoID)
	return nil
}

// Uploads data directly to targetFilePath in repoID
//
//   - repoID: name of the organisation repository
//   - data: raw data (string)
//   - targetFilePath: file to be created/updated on remote
//   - commitMessage: the message which will be added to the commit
//   - branch: the branch the data is to be pushed to
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func UploadRaw(repoID string, data string, targetFilePath string, commitMessage string, branch string) error {
	if err := uploadRaw(repoID, data, targetFilePath, commitMessage, branch); err != nil {
		return fmt.Errorf("could not upload raw data to repo %s: %w", repoID, err)
	}
	logger.Info.Printf("uploaded raw data to repo %s", repoID)
	return nil
}

// Adds a release to repoID with tagName & releaseName.
// Adds the path to the traces and the last graded timestamp to the release body.
//
//   - repoID: name of the organisation repository
//   - tagName: tag under which the release is to be created
//   - releaseName: name/title of the release
//   - tracesPath: path to the traces, used for adding a link to them to the release body
//   - graded: if set to true, the last graded timestamp in the release will be updated
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func NewRelease(repoID string, tagName string, releaseName string, tracesPath string, graded bool) error {
	if err := newRelease(repoID, tagName, releaseName, tracesPath, graded); err != nil {
		return err
	}
	logger.Info.Printf("added new release '%s' to %s", releaseName, repoID)
	return nil
}

// Gets the decoded file content from the specified repo, branch, and path as a string.
//
//   - repoID: name of the organisation repository
//   - branch: branch you would like to pull from
//   - filePath: file on the remote you would like to pull
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func GetDecodedFile(repoID string, branch string, filePath string) (content string, err error) {
	content, err = getFile(repoID, branch, filePath)
	if err != nil {
		return "", err
	}
	return content, nil
}
