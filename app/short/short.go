package short

import (
	"fmt"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/dao"
	"github.com/42-Short/shortinette/git"
	"github.com/42-Short/shortinette/logger"
)

type Short struct {
	Participants []dao.Participant
	Config       config.Config
	GitHubClient git.GithubService
}

func NewShort(participants []dao.Participant, config config.Config) (sh Short) {
	return Short{
		Participants: participants,
		Config:       config,
		GitHubClient: *git.NewGithubService(config.TokenGithub, config.OrgaGithub, config.BasePath),
	}
}

func (sh *Short) launchModule(moduleNumber int) (err error) {
	templateName, err := sh.GitHubClient.CreateModuleTemplate(moduleNumber)
	if err != nil {
		return fmt.Errorf("could not create template for module %02d: %v", moduleNumber, err)
	}

	for _, participant := range sh.Participants {
		repoName := fmt.Sprintf("%s-%02d", participant.IntraLogin, moduleNumber)
		description := fmt.Sprintf("Commit on the main branch with 'grademe' as a commit message to get graded. Minimum passing grade: %d", sh.Config.Modules[moduleNumber].MinimumScore)

		if err := sh.GitHubClient.NewRepo(templateName, repoName, true, description); err != nil {
			return fmt.Errorf("could not create new repo %s: %v", repoName, err)
		}

		if err := sh.GitHubClient.AddCollaborator(repoName, participant.GitHubLogin, "write"); err != nil {
			return fmt.Errorf("could not give %s write access to %s: %v", participant.GitHubLogin, repoName, err)
		}
	}

	return nil
}

func (sh *Short) Launch() (err error) {
	logger.Info.Println("launching Short now!")

	sh.launchModule(00)

	return nil
}
