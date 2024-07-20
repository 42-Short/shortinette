package git

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
)

// Clone a GitHub repo into targetDirectory.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func Clone(repoURL string, targetDirectory string) error {
	if err := get(repoURL, targetDirectory); err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}

// Check if repo exists, if not create it under the configured organisation.
// Also adds a webhook for easy recording of repository activity.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func Create(name string) error {
	if err := create(name); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not create repo: %w", err)
	}
	return nil
}

// Add a collaborator with the specified permissions to the repo.
//
//   - repoId: name of the organisation repository
//   - username: GitHub username of the collaborator
//   - permission: access level to be given to the user
//
// NOTE: Using this function will overwrite the user's previous rights - use test
// accounts, or you might lock yourself out of your repos.
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func AddCollaborator(repoId string, username string, permission string) error {
	if err := addCollaborator(repoId, username, permission); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not add %s to repo %s: %w", username, repoId, err)
	}
	return nil
}

// Add/Update a file on a repository
//
//   - repoId: name of the organisation repository
//   - localFilePath: source file to be uploaded
//   - targetFilePath: file to be created/updated on the remote
//
// See https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on GitHub configuration.
func UploadFile(repoId string, localFilePath string, targetFilePath string, commitMessage string, branch string) error {
	if err := uploadFile(repoId, localFilePath, targetFilePath, commitMessage, branch); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not upload %s to repo %s: %w", localFilePath, repoId, err)
	}
	logger.Info.Printf("%s successfully uploaded to %s/%s", localFilePath, repoId, targetFilePath)
	return nil
}

func NewRelease(repoId string, tagName string, releaseName string, formatReleaseName bool) error {
	if err := newRelease(repoId, tagName, releaseName, formatReleaseName); err != nil {
		return err
	}
	logger.Info.Printf("successfully added new release to %s", repoId)
	return nil
}

func getWaitTime(releaseName string) time.Duration {
	if releaseName == "" {
		logger.Info.Println("exercise has not been graded yet")
		return 0
	}
	waitTime, err := extractNumberFromString(releaseName)
	if err != nil {
		logger.Error.Printf("could not parse wait time from '%s'", releaseName)
		return 15 * time.Minute
	} else {
		return time.Duration(waitTime) * time.Minute
	}
}

func IsReadyToGrade(repoid string) (waitTime time.Duration, score int) {
	_, releaseName, releaseBody, err := getLatestRelease(repoid)
	if err != nil {
		logger.Error.Printf("failed getting the latest release for repo %s: %v", repoid, err)
		return 15 * time.Minute, 0
	}

	waitTime = getWaitTime(releaseName)

	if releaseBody == "" {
		return waitTime, 0
	}

	nameParts := strings.Split(releaseName, "/")
	if len(nameParts) == 0 {
		logger.Error.Printf("invalid release name format: %s", releaseName)
		return 15 * time.Minute, 0
	}
	score, err = strconv.Atoi(nameParts[0])
	if err != nil {
		score = 0
	}

	const timeStringLayout = "2006-01-02 15:04:05.999999999 -0700 MST"
	startIndex := len("last grading time: ")
	endIndex := strings.Index(releaseBody, "CEST") + len("CEST")
	lastGradingTimeStr := releaseBody[startIndex:endIndex]
	lastGradingTime, err := time.Parse(timeStringLayout, lastGradingTimeStr)
	if err != nil {
		logger.Error.Println("Error parsing last grading time:", err)
		return waitTime, score
	}

	timePassed := time.Since(lastGradingTime)

	if timePassed < waitTime {
		return waitTime - timePassed, score
	}

	return 0, score
}

func DeleteRepo(repoId string) error {
	if err := deleteRepo(repoId); err != nil {
		logger.Error.Println(err)
		return fmt.Errorf("could not delete repo %s: %w", repoId, err)
	}
	logger.Info.Printf("successfully deleted repo %s", repoId)
	return nil
}

// Actual implementation of the delete repo logic
func deleteRepo(repoId string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", os.Getenv("GITHUB_ORGANISATION"), repoId)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete repo: %s", resp.Status)
	}

	return nil
}
