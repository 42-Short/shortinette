package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func buildReleaseURL(repoID string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", os.Getenv("GITHUB_ORGANISATION"), repoID)
}

func buildLatestReleaseURL(repoID string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", os.Getenv("GITHUB_ORGANISATION"), repoID)
}

func createReleaseRequest(url string, token string, tagName string, releaseName string, body string) (*http.Request, error) {
	releaseDetails := map[string]interface{}{
		"tag_name": tagName,
		"name":     releaseName,
		"body":     body,
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

func createRelease(repo string, tagName string, releaseName string, body string) error {
	url := buildReleaseURL(repo)

	request, err := createReleaseRequest(url, os.Getenv("GITHUB_TOKEN"), tagName, releaseName, body)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

func newRelease(repoID string, tagName string, releaseName string, tracesPath string, graded bool) error {
	existingReleaseID, _, existingReleaseBody, err := getLatestRelease(repoID)
	if err != nil {
		return fmt.Errorf("could not check for existing release: %w", err)
	}

	if existingReleaseID != "" {
		if err := deleteRelease(repoID, existingReleaseID); err != nil {
			return fmt.Errorf("could not delete existing release: %w", err)
		}
	}

	newBody := existingReleaseBody
	if graded {
		newBody = fmt.Sprintf("**Last Graded:**\n- %s\n\n**Last Traces:**\n- https://github.com/42-Short/%s/tree/traces/%s", time.Now().Format("Monday, January 2, 2006 at 3:04 PM"), repoID, tracesPath)
	}

	if err := createRelease(repoID, tagName, releaseName, newBody); err != nil {
		return fmt.Errorf("could not create release: %w", err)
	}
	return nil
}

func getLatestRelease(repoID string) (id string, name string, body string, err error) {
	url := buildLatestReleaseURL(repoID)
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
	if err = json.NewDecoder(response.Body).Decode(&release); err != nil {
		return "", "", "", err
	}

	return fmt.Sprintf("%.0f", release["id"].(float64)), fmt.Sprintf("%s", release["name"]), fmt.Sprintf("%s", release["body"]), nil
}

func deleteRelease(repoID string, releaseID string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s", os.Getenv("GITHUB_ORGANISATION"), repoID, releaseID)
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
