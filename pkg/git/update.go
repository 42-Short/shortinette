package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getToken() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	return token, nil
}

func buildCollaboratorURL(repo, username string) string {
	return fmt.Sprintf("https://api.github.com/repos/42-Short/%s/collaborators/%s", repo, username)
}

func createCollaboratorRequest(url, token, permission string) (*http.Request, error) {
	collaboratorDetails := map[string]string{
		"permission": permission,
	}

	collaboratorDetailsJSON, err := json.Marshal(collaboratorDetails)
	if err != nil {
		return nil, fmt.Errorf("could not marshal collaborator details: %w", err)
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(collaboratorDetailsJSON))
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
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
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to add collaborator: %s", response.Status)
	}

	fmt.Println("Collaborator added successfully.")
	return nil
}

func addCollaborator(repo, username, permission string) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	url := buildCollaboratorURL(repo, username)

	request, err := createCollaboratorRequest(url, token, permission)
	if err != nil {
		return err
	}

	return sendRequest(request)
}
