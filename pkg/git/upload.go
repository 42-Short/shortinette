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

func buildPushURL(repo string, targetFilePath string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", os.Getenv("GITHUB_ORGANISATION"), repo, targetFilePath)
}

func getFileSHA(url, token string) (string, error) {
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

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get file SHA: %s, %s", resp.Status, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if sha, ok := result["sha"].(string); ok {
		return sha, nil
	}

	return "", fmt.Errorf("SHA not found in response")
}

func createPushRequest(url string, token string, targetFilePath string, commitMessage string, encodedContent string, sha string, branch string) (*http.Request, error) {
	requestDetails := map[string]interface{}{
		"message": commitMessage,
		"committer": map[string]string{
			"name":  os.Getenv("GITHUB_ADMIN"),
			"email": os.Getenv("GITHUB_EMAIL"),
		},
		"content": encodedContent,
		"path":    targetFilePath,
	}
	if sha != "" {
		requestDetails["sha"] = sha
	}

	if branch != "" {
		requestDetails["branch"] = branch
		logger.Info.Printf("pushing to branch: %s", branch)
	}

	requestDetailsJSON, err := json.Marshal(requestDetails)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestDetailsJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func uploadRaw(repoID string, data string, targetFilePath string, commitMessage string, branch string) (err error) {
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))

	url := buildPushURL(repoID, targetFilePath)
	shaURL := buildFileURL(repoID, branch, targetFilePath)

	sha, err := getFileSHA(shaURL, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}

	request, err := createPushRequest(url, os.Getenv("GITHUB_TOKEN"), targetFilePath, commitMessage, encodedData, sha, branch)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

func uploadFile(repoID string, localFilePath string, targetFilePath string, commitMessage string, branch string) error {
	originalFile, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("could not open original file: %w", err)
	}
	defer originalFile.Close()
	fileContent, err := io.ReadAll(originalFile)
	if err != nil {
		return err
	}
	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	url := buildPushURL(repoID, targetFilePath)

	sha, err := getFileSHA(url, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}

	request, err := createPushRequest(url, os.Getenv("GITHUB_TOKEN"), targetFilePath, commitMessage, encodedContent, sha, branch)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}
