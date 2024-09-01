// Package git provides functions for interacting with GitHub repositories, including
// cloning repositories, adding collaborators, uploading files, and creating releases.
package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// buildReleaseURL constructs the GitHub API URL for managing releases in a specific repository.
//
//   - repoID: the name of the repository
//
// Returns the URL as a string.
func buildReleaseURL(repoID string) (url string) {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", os.Getenv("GITHUB_ORGANISATION"), repoID)
}

// buildLatestReleaseURL constructs the GitHub API URL for retrieving the latest release of a specific repository.
//
//   - repoID: the name of the repository
//
// Returns the URL as a string.
func buildLatestReleaseURL(repoID string) (url string) {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", os.Getenv("GITHUB_ORGANISATION"), repoID)
}

// createReleaseRequest creates an HTTP request for creating a new GitHub release.
//
//   - url: the GitHub API URL for creating the release
//   - token: the GitHub authentication token
//   - tagName: the tag name for the release
//   - releaseName: the name/title of the release
//   - body: the body text of the release
//
// Returns the created HTTP request or an error if the request could not be created.
func createReleaseRequest(url string, token string, tagName string, releaseName string, body string) (request *http.Request, err error) {
	releaseDetails := map[string]interface{}{
		"tag_name": tagName,
		"name":     releaseName,
		"body":     body,
	}

	releaseDetailsJSON, err := json.Marshal(releaseDetails)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequest("POST", url, bytes.NewBuffer(releaseDetailsJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

// createRelease creates a new release in the specified GitHub repository.
//
//   - repo: the name of the repository
//   - tagName: the tag name for the release
//   - releaseName: the name/title of the release
//   - body: the body text of the release
//
// Returns an error if the release creation fails.
func createRelease(repo string, tagName string, releaseName string, body string) (err error) {
	url := buildReleaseURL(repo)

	request, err := createReleaseRequest(url, os.Getenv("GITHUB_TOKEN"), tagName, releaseName, body)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

// newRelease creates or updates a release for the specified repository, including handling
// the deletion of any existing release with the same tag.
//
//   - repoID: the name of the repository
//   - tagName: the tag name for the release
//   - releaseName: the name/title of the release
//   - tracesPath: the path to the traces, used for adding a link to them in the release body
//   - graded: if set to true, the last graded timestamp in the release will be updated
//
// Returns an error if the release creation or update fails.
func newRelease(repoID string, tagName string, releaseName string, tracesPath string, graded bool) (err error) {
	existingReleaseID, _, existingReleaseBody, err := getLatestRelease(repoID)
	if err != nil {
		return fmt.Errorf("could not check for existing release: %w", err)
	}

	if existingReleaseID != "" {
		if err := deleteRelease(repoID, existingReleaseID); err != nil {
			return fmt.Errorf("could not delete existing release: %w", err)
		}
	}

	newBody := existingReleaseBody
	if graded {
		newBody = fmt.Sprintf("**Last Graded:**\n- %s\n\n**Last Traces:**\n- https://github.com/%s/%s/tree/traces/%s", time.Now().Format("Monday, January 2, 2006 at 3:04 PM"), os.Getenv("GITHUB_ORGANISATION"), repoID, tracesPath)
	}

	if err := createRelease(repoID, tagName, releaseName, newBody); err != nil {
		return fmt.Errorf("could not create release: %w", err)
	}
	return nil
}

// getLatestRelease retrieves the ID, name, and body of the latest release in the specified repository.
//
//   - repoID: the name of the repository
//
// Returns the release ID, name, body, and an error if the retrieval fails.
func getLatestRelease(repoID string) (id string, name string, body string, err error) {
	url := buildLatestReleaseURL(repoID)
	token := os.Getenv("GITHUB_TOKEN")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return "", "", "", nil
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", "", "", fmt.Errorf("failed to get latest release with status %d: %s", response.StatusCode, string(body))
	}

	var release map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&release); err != nil {
		return "", "", "", err
	}

	return fmt.Sprintf("%.0f", release["id"].(float64)), fmt.Sprintf("%s", release["name"]), fmt.Sprintf("%s", release["body"]), nil
}

// deleteRelease deletes a specified release from a GitHub repository.
//
//   - repoID: the name of the repository
//   - releaseID: the ID of the release to be deleted
//
// Returns an error if the deletion fails.
func deleteRelease(repoID string, releaseID string) (err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s", os.Getenv("GITHUB_ORGANISATION"), repoID, releaseID)
	token := os.Getenv("GITHUB_TOKEN")

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to delete release with status %d: %s", response.StatusCode, string(body))
	}

	return nil
}
