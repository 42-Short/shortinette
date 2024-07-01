package endpoints

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/git"
	"github.com/42-Short/shortinette/internal/tester"
)

func CreateNewTeam(githubLogin string, projectId string) (err error) {
	repoName := fmt.Sprintf("%s-%s", githubLogin, projectId)
	if err = git.Create(repoName); err != nil {
		return err
	}
	if err = git.AddCollaborator(repoName, githubLogin, "push"); err != nil {
		return err
	}
	return nil
}

func TestSubmission(repoId string, testConfigPath string) (result map[string]error, err error) {
	result, err = tester.Run(testConfigPath, repoId, "studentcode")
	if err != nil {
		return nil, errors.NewInternalError(errors.ErrInternal, fmt.Sprintf("failed to run tests: %v", err))
	}
	return result, nil
}
