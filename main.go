package main

import (
	"errors"
	"fmt"
	"log"

	internalErrors "github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/git"
	"github.com/42-Short/shortinette/pkg/tester"
	"github.com/joho/godotenv"
)

func createNewTeam(githubLogin string, projectId string) (err error) {
	repoName := fmt.Sprintf("%s-%s", githubLogin, projectId)
	if err = git.Create(repoName); err != nil {
		return err
	}
	if err = git.AddCollaborator(repoName, githubLogin, "push"); err != nil {
		return err
	}
	return nil
}

func testSubmission(repoId string, testConfigPath string) (result string, err error) {
	if err := tester.Run(testConfigPath, repoId, "studentcode"); err != nil {
		var submissionErr *internalErrors.SubmissionError
		if errors.As(err, &submissionErr) {
			switch {
			case errors.Is(submissionErr.Err, internalErrors.ErrEmptyRepo):
				return "empty repository", nil
			case errors.Is(submissionErr.Err, internalErrors.ErrForbiddenItem):
				return "forbidden items used", nil
			case errors.Is(submissionErr.Err, internalErrors.ErrInvalidOutput):
				return "invalid output", nil
			case errors.Is(submissionErr.Err, internalErrors.ErrInvalidCompilation):
				return "invalid compilation", nil
			case errors.Is(submissionErr.Err, internalErrors.ErrRuntime):
				return "runtime error", nil
			default:
				return "", err
			}
		} else {
			return "", err
		}
	}
	return "success", nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}
	if err := createNewTeam("shortinette-test", "R00"); err != nil {
		log.Fatalf("could not create team: %s", err)
	}
	if result, err := testSubmission("shortinette-test-R00", "testconfig/R00.yaml"); err != nil {
		log.Fatalf("could not run tests: %s", err)
	} else {
		fmt.Printf("tests run successfully: %s", result)
	}
}
