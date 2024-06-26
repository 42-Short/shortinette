package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"os"
)

func RepoExists(repo string) (bool, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return false, fmt.Errorf("error: GITHUB_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://api.github.com/repos/42-Short/%s", repo)

	request, err := nethttp.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating HTTP request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &nethttp.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == nethttp.StatusOK {
		return true, nil
	} else if response.StatusCode == nethttp.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check if repository exists: %s", response.Status)
	}
}


func create(name string) error {
	exists, err := RepoExists(name)
	if err != nil {
		return fmt.Errorf("could not verify repository existence")
	}

	if exists {
		fmt.Printf("Repository %s already exists. Skipping creation.\n", name)
		return nil
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	url := "https://api.github.com/orgs/42-Short/repos"

	repoDetails := map[string]interface{}{
		"name":    name,
		"private": true,
	}
	repoDetailsJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return fmt.Errorf("could not marshal repository details: %w", err)
	}

	request, err := nethttp.NewRequest("POST", url, bytes.NewBuffer(repoDetailsJSON))
	if err != nil {
		return fmt.Errorf("could not create HTTP request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	client := &nethttp.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != nethttp.StatusCreated {
		return fmt.Errorf("could not create repository: %s", response.Status)
	}

	fmt.Println("Repository created successfully.")

	return nil
}
