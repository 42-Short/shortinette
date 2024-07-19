package git

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
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

func sendRequest(request *http.Request) error {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("request failed: %s, %s", response.Status, body)
	}
	logger.Info.Println("operation successful")
	return nil
}

func addCollaborator(repo, username, permission string) error {
	url := buildCollaboratorURL(repo, username)

	request, err := createCollaboratorRequest(url, os.Getenv("GITHUB_TOKEN"), permission)
	if err != nil {
		return err
	}

	return sendRequest(request)
}

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

func uploadFile(repoId string, localFilePath string, targetFilePath string, commitMessage string, branch string) error {
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

	url := buildPushURL(repoId, targetFilePath)

	sha, err := getFileSHA(url, os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}

	request, err := createPushRequest(url, os.Getenv("GITHUB_TOKEN"), targetFilePath, commitMessage, encodedContent, sha, branch)
	if err != nil {
		return err
	}

	return sendRequest(request)
}

func buildReleaseURL(repoId string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", os.Getenv("GITHUB_ORGANISATION"), repoId)
}

func buildLatestReleaseURL(repoId string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", os.Getenv("GITHUB_ORGANISATION"), repoId)
}

func createReleaseRequest(url string, token string, tagName string, releaseName string, body string, draft bool, prerelease bool) (*http.Request, error) {
	releaseDetails := map[string]interface{}{
		"tag_name":   tagName,
		"name":       releaseName,
		"body":       body,
		"draft":      draft,
		"prerelease": prerelease,
	}

	releaseDetailsJSON, err := json.Marshal(releaseDetails)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(releaseDetailsJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func createRelease(repo string, tagName string, releaseName string, body string, draft bool, prerelease bool) error {
	url := buildReleaseURL(repo)

	request, err := createReleaseRequest(url, os.Getenv("GITHUB_TOKEN"), tagName, releaseName, body, draft, prerelease)
	if err != nil {
		return err
	}

	return sendRequest(request)
}

func extractNumberFromString(s string) (int, error) {
	re := regexp.MustCompile(`(\d+)m`)
	matches := re.FindStringSubmatch(s)

	if len(matches) < 2 {
		return 0, fmt.Errorf("no number found")
	}

	number, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %v", err)
	}

	return number, nil
}

func newRelease(repoId string, tagName string, releaseName string, draft bool, prerelease bool, formatReleaseName bool) error {
	newWaitTime, newScore, currentScore := 0, 0, 0
	existingReleaseID, releaseTitle, body, err := getLatestRelease(repoId)
	if err != nil {
		return fmt.Errorf("could not check for existing release: %w", err)
	}
	if existingReleaseID != "" {
		newScore, err = strconv.Atoi(strings.Split(releaseName, "/")[0])
		if err != nil {
			return err
		}
		currentScore, err = strconv.Atoi(strings.Split(releaseTitle, "/")[0])
		if err != nil {
			return err
		}
		if err := deleteRelease(repoId, existingReleaseID); err != nil {
			return fmt.Errorf("could not delete existing release: %w", err)
		}
	}
	if newScore > currentScore || existingReleaseID == "" {
		newWaitTime = 15
	} else {
		oldWaitTime, err := extractNumberFromString(releaseTitle)
		if err != nil {
			return err
		}
		newWaitTime = min(oldWaitTime+15, 60)
	}
	newBody := ""
	if formatReleaseName {
		newBody = fmt.Sprintf("last grading time: %s", time.Now().String())
	} else {
		newBody = body
	}

	if formatReleaseName {
		releaseName = fmt.Sprintf("%s - retry in %dm", releaseName, newWaitTime)
	}

	if err := createRelease(repoId, tagName, releaseName, newBody, draft, prerelease); err != nil {
		return fmt.Errorf("could not create release: %w", err)
	}
	return nil
}

func getLatestRelease(repoId string) (string, string, string, error) {
	url := buildLatestReleaseURL(repoId)
	token := os.Getenv("GITHUB_TOKEN")

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return "", "", "", nil
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", "", "", fmt.Errorf("failed to get latest release with status %d: %s", response.StatusCode, string(body))
	}

	var release map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return "", "", "", err
	}

	return fmt.Sprintf("%.0f", release["id"].(float64)), fmt.Sprintf("%s", release["name"]), fmt.Sprintf("%s", release["body"]), nil
}

func deleteRelease(repoId string, releaseID string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s", os.Getenv("GITHUB_ORGANISATION"), repoId, releaseID)
	token := os.Getenv("GITHUB_TOKEN")

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to delete release with status %d: %s", response.StatusCode, string(body))
	}

	return nil
}
