package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"os"
)

func addCollaborator(repo string, username string, permission string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://api.github.com/repos/42-Short/%s/collaborators/%s", repo, username)

	collaboratorDetails := map[string]string{
		"permission": permission,
	}

	collaboratorDetailsJSON, err := json.Marshal(collaboratorDetails)
	if err != nil {
		return fmt.Errorf("could not marshal collaborator details")
	}

	request, err := nethttp.NewRequest("PUT", url, bytes.NewBuffer(collaboratorDetailsJSON))
	if err != nil {
		return fmt.Errorf("could not create HTTP request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	client := &nethttp.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != nethttp.StatusCreated && response.StatusCode != nethttp.StatusNoContent {
		return fmt.Errorf("failed to add collaborator: %s", response.Status)
	}

	fmt.Println("Collaborator added successfully.")
	return nil
}
