package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"bytes"
	"encoding/base64"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/joho/godotenv"
)

// buildPushURL constructs the GitHub API URL for pushing a file to a specific repository
// and target file path.
//
//   - repo: the name of the repository
//   - targetFilePath: the path of the file in the repository
//
// Returns the URL as a string.
func buildPushURL(repo string, targetFilePath string) (url string) {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", "Short-Test-Orga", repo, targetFilePath)
}

func sendHTTPRequest(request *http.Request) (response *http.Response, err error) {
	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(response.Body)
		logger.Info.Printf("x-ratelimit-remaining: %s\n", response.Header.Get("x-ratelimit-remaining"))
		logger.Info.Printf("x-ratelimit-reset: %s\n", response.Header.Get("x-ratelimit-reset"))
		logger.Info.Printf("retry-after: %s\n", response.Header.Get("retry-after"))
		logger.Info.Printf("response body: %s", body)
		return response, fmt.Errorf("request failed: %s, %s", response.Status, body)
	}
	return response, nil
}

// getFileSHA retrieves the SHA of a file in a specific repository at the given URL.
//
//   - url: the GitHub API URL to retrieve the file SHA
//   - token: the GitHub authentication token
//
// Returns the file's SHA as a string and an error if the operation fails.
func getFileSHA(url, token string) (sha string, err error) {
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

// createPushRequest creates an HTTP request for pushing a file to a GitHub repository.
//
//   - url: the GitHub API URL for pushing the file
//   - token: the GitHub authentication token
//   - targetFilePath: the path of the file in the repository
//   - commitMessage: the commit message to use
//   - encodedContent: the base64-encoded content of the file
//   - sha: the SHA of the file being updated (optional, can be empty for new files)
//   - branch: the branch to push to (optional)
//
// Returns the created HTTP request or an error if the request could not be created.
func createPushRequest(url string, token string, targetFilePath string, commitMessage string, encodedContent string, sha string, branch string) (request *http.Request, err error) {
	requestDetails := map[string]interface{}{
		"message": commitMessage,
		"committer": map[string]string{
			"name":  "winstonallo",
			"email": "arthurbiedchar@gmail.com",
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

	request, err = http.NewRequest("PUT", url, bytes.NewBuffer(requestDetailsJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func buildFileURL(repoID string, branch string, filePath string) (url string) {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", "Short-Test-Orga", repoID, filePath, branch)
}

// uploadRaw uploads raw data directly to a GitHub repository at the specified file path.
//
//   - repoID: the name of the repository
//   - data: the raw data to be uploaded as a string
//   - targetFilePath: the path of the file in the repository
//   - commitMessage: the commit message to use
//   - branch: the branch to push to (optional)
//
// Returns an error if the upload process fails.
func uploadRaw(repoID string, data string, targetFilePath string, commitMessage string, branch string) (err error) {
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))

	url := buildPushURL(repoID, targetFilePath)
	shaURL := buildFileURL(repoID, branch, targetFilePath)

	sha, err := getFileSHA(shaURL, (os.Getenv("GITHUB_TOKEN")))
	if err != nil {
		return err
	}

	request, err := createPushRequest(url, (os.Getenv("GITHUB_TOKEN")), targetFilePath, commitMessage, encodedData, sha, branch)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

// uploadFile uploads a local file to a GitHub repository at the specified file path.
//
//   - repoID: the name of the repository
//   - localFilePath: the path of the local file to be uploaded
//   - targetFilePath: the path of the file in the repository
//   - commitMessage: the commit message to use
//   - branch: the branch to push to (optional)
//
// Returns an error if the upload process fails.
func uploadFile(repoID string, localFilePath string, targetFilePath string, commitMessage string, branch string) (err error) {
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

	sha, err := getFileSHA(url, (os.Getenv("GITHUB_TOKEN")))
	if err != nil {
		return err
	}

	request, err := createPushRequest(url, (os.Getenv("GITHUB_TOKEN")), targetFilePath, commitMessage, encodedContent, sha, branch)
	if err != nil {
		return err
	}

	if _, err := sendHTTPRequest(request); err != nil {
		return err
	}
	return nil
}

type Participant struct {
	IntraLogin string `json:"intra_login"`
}

type ShortConfig struct {
	Participants []Participant `json:"participants"`
}

func deleteRepo(repoID string) bool {
	org := os.Getenv("GITHUB_ORGANISATION")
	token := (os.Getenv("GITHUB_TOKEN"))
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", org, repoID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		log.Printf("Successfully deleted repo %s", repoID)
		return true
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to delete repo %s: %d %s", repoID, resp.StatusCode, string(bodyBytes))
		return false
	}
}

func main() {
	logger.InitializeStandardLoggers("TEST")
	godotenv.Load()

	if len(os.Args) != 4 {
		fmt.Println("usage: go run . <module> <path-to-new-subject> <commit-message>")
		return
	}

	moduleName := os.Args[1]
	subjectPath := os.Args[2]
	commitMessage := os.Args[3]

	config, _ := os.ReadFile("config/final-participants.json")

	var shortConfig ShortConfig
	if err := json.Unmarshal(config, &shortConfig); err != nil {
		log.Fatalf("Error parsing %s: %v", os.Getenv("CONFIG_PATH"), err)
	}

	for _, participant := range shortConfig.Participants {
		if err := uploadFile(fmt.Sprintf("%s-%s", participant.IntraLogin, moduleName), subjectPath, "README.md", commitMessage, "main"); err != nil {
			fmt.Println(err)
			fmt.Println((os.Getenv("GITHUB_TOKEN")))
		} else {
			fmt.Printf("updated repo for %s\n", participant.IntraLogin)
		}
	}
}
