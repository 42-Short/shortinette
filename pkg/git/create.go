package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// buildRepoURL constructs the GitHub API URL for the repository
func buildRepoURL(repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/42-Short/%s", repo)
}

// buildCreateRepoURL constructs the GitHub API URL for creating a repository in the organization
func buildCreateRepoURL() string {
	return "https://api.github.com/orgs/42-Short/repos"
}

// createHTTPRequest creates a new HTTP request with the provided method, URL, token, and body
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

// sendHTTPRequest sends the provided HTTP request and returns the response
func sendHTTPRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	return response, nil
}

// RepoExists checks if the specified repository exists in the organization
func RepoExists(repo string) (bool, error) {
	token, err := getToken()
	if err != nil {
		return false, err
	}

	url := buildRepoURL(repo)
	request, err := createHTTPRequest("GET", url, token, nil)
	if err != nil {
		return false, err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return true, nil
	} else if response.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check if repository exists: %s", response.Status)
	}
}

// createRepository creates a new repository in the organization with the specified name
func createRepository(name string) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	url := buildCreateRepoURL()
	repoDetails := map[string]interface{}{
		"name":    name,
		"private": true,
	}
	repoDetailsJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return fmt.Errorf("could not marshal repository details: %w", err)
	}

	request, err := createHTTPRequest("POST", url, token, repoDetailsJSON)
	if err != nil {
		return err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create repository: %s", response.Status)
	}

	fmt.Println("Repository created successfully.")
	return nil
}

func create(name string) error {
	exists, err := RepoExists(name)
	if err != nil {
		return fmt.Errorf("could not verify repository existence: %w", err)
	}

	if exists {
		fmt.Printf("Repository %s already exists. Skipping creation.\n", name)
		return nil
	}

	return createRepository(name)
}
