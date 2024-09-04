// git provides functions for interacting with GitHub repositories, including
// cloning repositories, adding collaborators, uploading files, and creating releases.
package git

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// buildCollaboratorURL constructs the GitHub API URL for managing a collaborator in a specific repository.
//
//   - repo: the name of the repository
//   - username: the GitHub username of the collaborator
//
// Returns the URL as a string.
func buildCollaboratorURL(repo string, username string) (url string) {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators/%s", os.Getenv("GITHUB_ORGANISATION"), repo, username)
}

// createCollaboratorRequest creates an HTTP request for adding or updating a collaborator
// in a GitHub repository with a specific permission level.
//
//   - url: the GitHub API URL for adding the collaborator
//   - token: the GitHub authentication token
//   - permission: the permission level to be granted to the collaborator (e.g., "push", "pull")
//
// Returns the created HTTP request or an error if the request could not be created.
func createCollaboratorRequest(url string, token string, permission string) (request *http.Request, err error) {
	collaboratorDetails := map[string]string{
		"permission": permission,
	}

	collaboratorDetailsJSON, err := json.Marshal(collaboratorDetails)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequest("PUT", url, bytes.NewBuffer(collaboratorDetailsJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// addCollaborator adds a collaborator to a GitHub repository with the specified permission level.
//
//   - repo: the name of the repository
//   - username: the GitHub username of the collaborator
//   - permission: the permission level to be granted to the collaborator (e.g., "push", "pull")
//
// Returns an error if the operation fails.
func addCollaborator(repoID string, username string, permission string) (err error) {
	url := buildCollaboratorURL(repoID, username)

	request, err := createCollaboratorRequest(url, os.Getenv("GITHUB_TOKEN"), permission)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

// buildFileURL constructs the GitHub API URL for retrieving a file from a specific repository
// branch and file path.
//
//   - repoID: the name of the repository
//   - branch: the branch from which to retrieve the file
//   - filePath: the path of the file in the repository
//
// Returns the URL as a string.
func buildFileURL(repoID string, branch string, filePath string) (url string) {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", os.Getenv("GITHUB_ORGANISATION"), repoID, filePath, branch)
}

// getFile retrieves the content of a file from a specific branch of a GitHub repository.
//
//   - repoID: the name of the repository
//   - branch: the branch from which to retrieve the file
//   - filePath: the path of the file in the repository
//
// Returns the file content as a string and an error if the operation fails.
func getFile(repoID string, branch string, filePath string) (content string, err error) {
	url := buildFileURL(repoID, branch, filePath)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_TOKEN")))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("failed to get %s from %s: %d - %s", filePath, repoID, response.StatusCode, string(body))
	}

	var contentJSON struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(response.Body).Decode(&content); err != nil {
		return "", err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(contentJSON.Content)
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}
