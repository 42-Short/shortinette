package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/42-Short/shortinette/internal/datastructures"
)

func buildRepoURL(repo string, env datastructures.Environment) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", env.Organisation, repo)
}

func buildCreateRepoURL(env datastructures.Environment) string {
	return fmt.Sprintf("https://api.github.com/orgs/%s/repos", env.Organisation)
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

func sendHTTPRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	return response, nil
}

func RepoExists(repo string, env datastructures.Environment) (bool, error) {
	url := buildRepoURL(repo, env)
	request, err := createHTTPRequest("GET", url, env.Token, nil)
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

func createRepository(name string, env datastructures.Environment) error {
	url := buildCreateRepoURL(env)
	repoDetails := map[string]interface{}{
		"name":    name,
		"private": true,
	}
	repoDetailsJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return fmt.Errorf("could not marshal repository details: %w", err)
	}

	request, err := createHTTPRequest("POST", url, env.Token, repoDetailsJSON)
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

func create(name string, env datastructures.Environment) error {
	exists, err := RepoExists(name, env)
	if err != nil {
		return fmt.Errorf("could not verify repository existence: %w", err)
	}

	if exists {
		fmt.Printf("Repository %s already exists. Skipping creation.\n", name)
		return nil
	}

	return createRepository(name, env)
}
