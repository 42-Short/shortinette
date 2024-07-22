package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getDefaultBranchSHA(repo, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/main", os.Getenv("GITHUB_ORGANISATION"), repo)
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

func createBranch(repo, token, branchName, sha string) error {
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
