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
		return fmt.Errorf(response.Status)
	}

	logger.Info.Printf("added webhook to %s", repoID)
	return nil
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

func createRepository(name string) (err error) {
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
		return fmt.Errorf(response.Status)
	}

	if err := addWebhook(name); err != nil {
		return fmt.Errorf("failed to add webhook: %w", err)
	}

	logger.Info.Printf("repository %s created in %s", name, os.Getenv("GITHUB_ORGANISATION"))
	return nil
}

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
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to make initial commit: %s, %s", resp.Status, body)
	}

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

	if err := createRepository(name); err != nil {
		return err
	}

	if err := initialCommit(name, os.Getenv("GITHUB_TOKEN")); err != nil {
		return err
	}

	sha, err := getDefaultBranchSHA(name, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}
	if err := createBranch(name, os.Getenv("GITHUB_TOKEN"), "traces", sha); err != nil {
		return err
	}
	// adding branch protections to private organisation repos only works with github enterprise
	// if err := addBranchProtection(name, "traces"); err != nil {
	// 	return err
	// }
	return nil
}
