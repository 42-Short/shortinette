package git

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func buildCollaboratorURL(repo string, username string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators/%s", os.Getenv("GITHUB_ORGANISATION"), repo, username)
}

func createCollaboratorRequest(url string, token string, permission string) (*http.Request, error) {
	collaboratorDetails := map[string]string{
		"permission": permission,
	}

	collaboratorDetailsJSON, err := json.Marshal(collaboratorDetails)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(collaboratorDetailsJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func addCollaborator(repo, username, permission string) (err error) {
	url := buildCollaboratorURL(repo, username)

	request, err := createCollaboratorRequest(url, os.Getenv("GITHUB_TOKEN"), permission)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

func buildFileURL(repoID, branch, filePath string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", os.Getenv("GITHUB_ORGANISATION"), repoID, filePath, branch)
}

func getFile(repoID, branch, filePath string) (string, error) {
	url := buildFileURL(repoID, branch, filePath)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GITHUB_TOKEN")))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("failed to get %s from %s: %d - %s", filePath, repoID, response.StatusCode, string(body))
	}

	var content struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(response.Body).Decode(&content); err != nil {
		return "", err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}
