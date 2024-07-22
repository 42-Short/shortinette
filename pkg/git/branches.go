package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/42-Short/shortinette/internal/logger"
)

func getDefaultBranchSHA(repoID string, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/main", os.Getenv("GITHUB_ORGANISATION"), repoID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get default branch SHA: %s, %s", resp.Status, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if sha, ok := result["object"].(map[string]interface{})["sha"].(string); ok {
		return sha, nil
	}

	return "", fmt.Errorf("SHA not found in response")
}

func createBranch(repo string, token string, branchName string, sha string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", os.Getenv("GITHUB_ORGANISATION"), repo)
	requestBody := map[string]interface{}{
		"ref": fmt.Sprintf("refs/heads/%s", branchName),
		"sha": sha,
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusUnprocessableEntity {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create branch: %s, %s", resp.Status, body)
	}

	return nil
}

func addBranchProtection(repoID string, branch string) (err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s/protection", os.Getenv("GITHUB_ORGANISATION"), repoID, branch)
	requestBody := map[string]interface{}{
		"restrictons": map[string]interface{}{
			"users": []string{os.Getenv("GITHUB_ADMIN")},
		},
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := createHTTPRequest("PUT", url, os.Getenv("GITHUB_TOKEN"), requestBodyJSON)
	if err != nil {
		return err
	}

	response, err := sendHTTPRequest(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}
	logger.Info.Printf("added protection to branch %s/%s", repoID, branch)
	return nil
}
