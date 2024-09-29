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
	"strconv"
	"time"

	"github.com/42-Short/shortinette/pkg/logger"
)

// addWebhook adds a webhook to the specified repository to listen for push events.
//
//   - repoID: the name of the repository to which the webhook will be added
//
// Returns an error if the webhook cannot be added.
func addWebhook(repoID string) (err error) {
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
func checkResponseStatus(response *http.Response) (ok bool, err error) {
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
func RepoExists(repo string) (exists bool, err error) {
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
func createRepository(name string, withWebhook bool) (callsRemaining int, reset time.Time, err error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", os.Getenv("GITHUB_ORGANISATION"))
	repoDetails := map[string]interface{}{
		"name":    name,
		"private": true,
	}
	repoDetailsJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return 0, time.Time{}, err
	}

	request, err := createHTTPRequest("POST", url, os.Getenv("GITHUB_TOKEN"), repoDetailsJSON)
	if err != nil {
		return 0, time.Time{}, err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return 0, time.Time{}, err
	}
	defer response.Body.Close()

	callsRemaining, err = strconv.Atoi(string(response.Header.Get("x-ratelimit-remaining")[0]))
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("x-ratelimit-remaining header could not be parsed: %v", err)
	}

	resetStr := response.Header.Get("x-ratelimit-reset")
	resetUnix, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		logger.Error.Printf("x-ratelimit-reset header could not be parsed: %v", err)
		reset = time.Now().Add(time.Hour)
	} else {
		reset = time.Unix(resetUnix, 0)
	}

	if response.StatusCode != http.StatusCreated {
		return 0, time.Time{}, fmt.Errorf("could not create repository %s: %s", name, response.Status)
	}

	if withWebhook {
		if err := addWebhook(name); err != nil {
			return 0, time.Time{}, fmt.Errorf("failed to add webhook: %w", err)
		}
	}

	logger.Info.Printf("repository %s created in %s", name, os.Getenv("GITHUB_ORGANISATION"))
	return callsRemaining, reset, err
}

// initialCommit creates an initial commit in the newly created repository with a README.md file.
//
//   - repo: the name of the repository where the commit will be made
//   - token: the GitHub authentication token
//
// Returns an error if the initial commit fails.
func initialCommit(repo, token string) (err error) {
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
//   - withWebhook: bool indicating whether a webhook should be added to the repo
//   - additionalBranche: variadic list of strings, representing all branches that should be created by default on the repo
//
// Returns an error if the repository creation, initial commit, or branch creation fails.
func create(name string, withWebhook bool, additionalBranches ...string) (callsRemaining int, reset time.Time, err error) {
	exists, err := RepoExists(name)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("could not verify repository existence: %w", err)
	}

	if exists {
		logger.Info.Printf("repository %s already exists. Skipping creation\n", name)
		return 0, time.Time{}, nil
	}

	if callsRemaining, reset, err = createRepository(name, withWebhook); err != nil {
		return 0, time.Time{}, nil
	}

	if err := initialCommit(name, os.Getenv("GITHUB_TOKEN")); err != nil {
		return 0, time.Time{}, nil
	}

	sha, err := getDefaultBranchSHA(name, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return 0, time.Time{}, nil
	}
	for _, branch := range additionalBranches {
		if err := createBranch(name, os.Getenv("GITHUB_TOKEN"), branch, sha); err != nil {
			return 0, time.Time{}, nil
		}
	}
	return callsRemaining, reset, nil
}
