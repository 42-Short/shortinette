// Package git provides functions for interacting with GitHub repositories, including
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

	"github.com/42-Short/shortinette/pkg/logger"
)

// addWebhook adds a webhook to the specified repository to listen for push events.
//
//   - repoID: the name of the repository to which the webhook will be added
//
// Returns an error if the webhook cannot be added.
func addWebhook(repoID string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", os.Getenv("GITHUB_ORGANISATION"), repoID)
	webhookConfig := map[string]interface{}{
		"name":   "web",
		"active": true,
		"events": []string{"push"},
		"config": map[string]string{
			"url":          os.Getenv("WEBHOOK_URL"),
			"content_type": "json",
		},
	}
	webhookConfigJSON, err := json.Marshal(webhookConfig)
	if err != nil {
		return err
	}

	request, err := createHTTPRequest("POST", url, os.Getenv("GITHUB_TOKEN"), webhookConfigJSON)
	if err != nil {
		return err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not add webhook to %s: %s", repoID, response.Status)
	}

	logger.Info.Printf("added webhook to %s", repoID)
	return nil
}

// checkResponseStatus checks the HTTP response status code and returns true if the response
// status code is 200 OK or an error if the status code indicates a failure.
//
//   - response: the HTTP response to check
//
// Returns a boolean indicating success or failure, and an error if the status code indicates a problem.
func checkResponseStatus(response *http.Response) (bool, error) {
	if response.StatusCode == http.StatusOK {
		return true, nil
	} else if response.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("error response: %s", response.Status)
	}
}

// RepoExists checks if a GitHub repository exists for the given repo name.
//
//   - repo: the name of the repository to check
//
// Returns a boolean indicating whether the repository exists and an error if the check fails.
func RepoExists(repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", os.Getenv("GITHUB_ORGANISATION"), repo)

	request, err := createHTTPRequest("GET", url, os.Getenv("GITHUB_TOKEN"), nil)
	if err != nil {
		return false, err
	}

	response, err := sendHTTPRequest(request)
	if err != nil && response.StatusCode != http.StatusNotFound {
		return false, err
	}
	defer response.Body.Close()

	return checkResponseStatus(response)
}

// createRepository creates a new GitHub repository with the specified name under the
// configured organization.
//
//   - name: the name of the repository to create
//
// Returns an error if the repository creation fails.
func createRepository(name string, withWebhook bool) (err error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", os.Getenv("GITHUB_ORGANISATION"))
	repoDetails := map[string]interface{}{
		"name":    name,
		"private": true,
	}
	repoDetailsJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return err
	}

	request, err := createHTTPRequest("POST", url, os.Getenv("GITHUB_TOKEN"), repoDetailsJSON)
	if err != nil {
		return err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create repository %s: %s", name, response.Status)
	}

	if withWebhook {
		if err := addWebhook(name); err != nil {
			return fmt.Errorf("failed to add webhook: %w", err)
		}
	}

	logger.Info.Printf("repository %s created in %s", name, os.Getenv("GITHUB_ORGANISATION"))
	return nil
}

// initialCommit creates an initial commit in the newly created repository with a README.md file.
//
//   - repo: the name of the repository where the commit will be made
//   - token: the GitHub authentication token
//
// Returns an error if the initial commit fails.
func initialCommit(repo, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/README.md", os.Getenv("GITHUB_ORGANISATION"), repo)
	requestBody := map[string]interface{}{
		"message": "Initial commit",
		"content": base64.StdEncoding.EncodeToString([]byte("")),
		"branch":  "main",
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to make initial commit: %s, %s", response.Status, body)
	}

	return nil
}

// create handles the process of creating a new repository if it does not already exist,
// making an initial commit, and creating necessary branches.
//
//   - name: the name of the repository to create
//
// Returns an error if the repository creation, initial commit, or branch creation fails.
func create(name string, withWebhook bool, additionalBranches ...string) error {
	exists, err := RepoExists(name)
	if err != nil {
		return fmt.Errorf("could not verify repository existence: %w", err)
	}

	if exists {
		logger.Info.Printf("repository %s already exists. Skipping creation\n", name)
		return nil
	}

	if err := createRepository(name, withWebhook); err != nil {
		return err
	}

	if err := initialCommit(name, os.Getenv("GITHUB_TOKEN")); err != nil {
		return err
	}

	sha, err := getDefaultBranchSHA(name, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}
	for _, branch := range additionalBranches {
		if err := createBranch(name, os.Getenv("GITHUB_TOKEN"), branch, sha); err != nil {
			return err
		}
	}
	return nil
}
