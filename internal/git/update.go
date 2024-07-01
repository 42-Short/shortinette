package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func buildCollaboratorURL(repo, username string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators/%s", os.Getenv("GITHUB_ORGANISATION"), repo, username)
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
	url := buildCollaboratorURL(repo, username)

	request, err := createCollaboratorRequest(url, os.Getenv("GITHUB_TOKEN"), permission)
	if err != nil {
		return err
	}

	return sendRequest(request)
}