package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/42-Short/shortinette/internal/logger"
)

func addWebhook(repoId string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", os.Getenv("GITHUB_ORGANISATION"), repoId)
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
		return fmt.Errorf(response.Status)
	}

	logger.Info.Printf("webhook added successfully to %s", repoId)
	return nil
}

func buildRepoURL(repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", os.Getenv("GITHUB_ORGANISATION"), repo)
}

func buildCreateRepoURL() string {
	return fmt.Sprintf("https://api.github.com/orgs/%s/repos", os.Getenv("GITHUB_ORGANISATION"))
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

func checkResponseStatus(response *http.Response) (bool, error) {
	if response.StatusCode == http.StatusOK {
		return true, nil
	} else if response.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf(response.Status)
	}
}

func RepoExists(repo string) (bool, error) {
	url := buildRepoURL(repo)
	request, err := createHTTPRequest("GET", url, os.Getenv("GITHUB_TOKEN"), nil)
	if err != nil {
		return false, err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	return checkResponseStatus(response)
}

func createRepository(name string) error {
	url := buildCreateRepoURL()
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
		return fmt.Errorf(response.Status)
	}

	if err := addWebhook(name); err != nil {
		return fmt.Errorf("failed to add webhook: %w", err)
	}

	fmt.Println("repository created successfully")
	return nil
}

func create(name string) error {
	exists, err := RepoExists(name)
	if err != nil {
		return fmt.Errorf("could not verify repository existence: %w", err)
	}

	if exists {
		logger.Info.Printf("repository %s already exists. Skipping creation\n", name)
		return nil
	}

	return createRepository(name)
}