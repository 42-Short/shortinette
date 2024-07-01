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

func testSubmission(repoId string, testConfigPath string) (result string) {
	if err := tester.Run(testConfigPath, repoId, "studentcode"); err != nil {
		var submissionErr *internalErrors.SubmissionError
		if errors.As(err, &submissionErr) {
			switch {
			case errors.Is(submissionErr.Err, internalErrors.ErrEmptyRepo):
				return "empty repository"
			case errors.Is(submissionErr.Err, internalErrors.ErrForbiddenItem):
				return "forbidden items used"
			case errors.Is(submissionErr.Err, internalErrors.ErrInvalidOutput):
				return "invalid output"
			case errors.Is(submissionErr.Err, internalErrors.ErrInvalidCompilation):
				return "invalid compilation"
			case errors.Is(submissionErr.Err, internalErrors.ErrRuntime):
				return "runtime error"
			default:
				return "unknown error"
			}
		} else {
			return "unknown error"
		}
	}
	return "success"
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}
	if err := createNewTeam("shortinette-test", "R00"); err != nil {
		log.Fatalf("could not create team: %s", err)
	}
}
