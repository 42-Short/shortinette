package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"os"
)

func doCreate(name string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	url := "https://api.github.com/orgs/42-Short/repos"

	repoDetails := map[string]interface{}{
		"name":    name,
		"private": true,
	}
	repoDetailsJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return fmt.Errorf("could not marshal repository details: %w", err)
	}

	request, err := nethttp.NewRequest("POST", url, bytes.NewBuffer(repoDetailsJSON))
	if err != nil {
		return fmt.Errorf("could not create HTTP request: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	client := &nethttp.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != nethttp.StatusCreated {
		return fmt.Errorf("could not create repository: %s", response.Status)
	}

	fmt.Println("Repository created successfully.")
	return nil
}
