// git provides functions for interacting with GitHub repositories, including
// cloning repositories, adding collaborators, uploading files, and creating releases.
package git

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/42-Short/shortinette/pkg/logger"
)

// sendHTTPRequest sends an HTTP request and checks the response status.
// It returns the response or an error if the request fails.
//
//   - request: the HTTP request to be sent
//
// Returns the HTTP response or an error if the request fails.
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

// createHTTPRequest creates an HTTP request with the specified method, URL, authorization
// token, and body.
//
//   - method: the HTTP method (e.g., "GET", "POST", "PUT", etc.)
//   - url: the URL for the request
//   - token: the authorization token to be included in the request header
//   - body: the request body as a byte slice
//
// Returns the created HTTP request or an error if the request could not be created.
func createHTTPRequest(method string, url string, token string, body []byte) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// Clone clones a GitHub repository from the specified repoURL into the targetDirectory.
//
//   - repoURL: the URL of the repository to clone
//   - targetDirectory: the directory where the repository should be cloned
//
// Returns an error if the cloning process fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func Clone(repoURL string, targetDirectory string) (err error) {
	if err = get(repoURL, targetDirectory); err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

// Create checks if a repository exists, and if not, creates it under the configured
// organization. It also adds a webhook for easy recording of repository activity.
//
//   - name: the name of the repository to create
//   - withWebhook: bool indicating whether to add a webhook to the repo, allowing the server to listen for changes
//   - additionalBranches: variadic list of branches you would like to be created
//
// Returns an error if the repository creation process fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func Create(name string, withWebhook bool, additionalBranches ...string) (callsRemaining int, reset time.Time, err error) {
	if callsRemaining, reset, err = create(name, withWebhook, additionalBranches...); err != nil {
		logger.Error.Println(err)
		return 0, time.Time{}, fmt.Errorf("could not create repo: %w", err)
	}
	return callsRemaining, reset, nil
}

// AddCollaborator adds a collaborator with the specified permissions to the repository.
//
//   - repoID: the name of the organization repository
//   - username: the GitHub username of the collaborator
//   - permission: the access level to be given to the user
//
// NOTE: Using this function will overwrite the user's previous rights. Use test
// accounts, or you might lock yourself out of your repos.
//
// Returns an error if the process of adding the collaborator fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func AddCollaborator(repoID string, username string, permission string) (callsRemaining int, reset time.Time, err error) {
	callsRemaining, reset, err = addCollaborator(repoID, username, permission)
	if err != nil {
		logger.Error.Println(err)
		return 0, time.Time{}, fmt.Errorf("could not add %s to repo %s: %w", username, repoID, err)
	}
	return callsRemaining, reset, nil
}

// UploadFile adds or updates a file on a repository.
//
//   - repoID: the name of the organization repository
//   - localFilePath: the source file to be uploaded
//   - targetFilePath: the file to be created/updated on the remote
//   - commitMessage: the message to be added to the commit
//   - branch: the branch to push the data to
//
// Returns an error if the file upload process fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func UploadFile(repoID string, localFilePath string, targetFilePath string, commitMessage string, branch string) (callsRemaining int, reset time.Time, err error) {
	if err := uploadFile(repoID, localFilePath, targetFilePath, commitMessage, branch); err != nil {
		return fmt.Errorf("could not upload %s to repo %s: %w", localFilePath, repoID, err)
	}
	logger.Info.Printf("uploaded %s to repo %s", localFilePath, repoID)
	return nil
}

// UploadRaw uploads raw data directly to the specified targetFilePath in the repository.
//
//   - repoID: the name of the organization repository
//   - data: the raw data to be uploaded as a string
//   - targetFilePath: the file to be created/updated on the remote
//   - commitMessage: the message to be added to the commit
//   - branch: the branch to push the data to
//
// Returns an error if the data upload process fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func UploadRaw(repoID string, data string, targetFilePath string, commitMessage string, branch string) (err error) {
	if err := uploadRaw(repoID, data, targetFilePath, commitMessage, branch); err != nil {
		return fmt.Errorf("could not upload raw data to repo %s: %w", repoID, err)
	}
	logger.Info.Printf("uploaded raw data to repo %s", repoID)
	return nil
}

// NewRelease adds a release to the specified repository with the provided tagName and
// releaseName. It also adds the path to the traces and the last graded timestamp to the
// release body if graded is set to true.
//
//   - repoID: the name of the organization repository
//   - tagName: the tag under which the release is to be created
//   - releaseName: the name/title of the release
//   - tracesPath: the path to the traces, used for adding a link to them in the release body
//   - graded: if set to true, the last graded timestamp in the release will be updated
//
// Returns an error if the release creation process fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func NewRelease(repoID string, tagName string, releaseName string, tracesPath string, graded bool) (err error) {
	if err := newRelease(repoID, tagName, releaseName, tracesPath, graded); err != nil {
		return err
	}
	logger.Info.Printf("added new release '%s' to %s", releaseName, repoID)
	return nil
}

// GetDecodedFile retrieves the decoded content of a file from the specified repository,
// branch, and file path as a string.
//
//   - repoID: the name of the organization repository
//   - branch: the branch to pull the file from
//   - filePath: the path to the file on the remote repository
//
// Returns the file content as a string and an error if the retrieval process fails.
//
// See https://github.com/42-Short/shortinette/README.md for details on GitHub configuration.
func GetDecodedFile(repoID string, branch string, filePath string) (content string, err error) {
	content, err = getFile(repoID, branch, filePath)
	if err != nil {
		return "", err
	}
	return content, nil
}
