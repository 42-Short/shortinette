package short

import (
	"fmt"
	"time"

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

func (sh *Short) getCurrentModuleIdx() (moduleIdx int) {
	for idx := range sh.Config.Modules {
		if time.Now().After(sh.Config.StartTime.Add(sh.Config.ModuleDuration*time.Duration(idx))) && time.Now().Before(sh.Config.StartTime.Add(sh.Config.ModuleDuration*time.Duration(idx+1))) {
			return idx
		}
	}

	return -1
}

func (sh *Short) schedule() {
	for time.Now().Before(sh.Config.StartTime) {
		time.Sleep(time.Until(sh.Config.StartTime))
	}

	currentModuleIdx := sh.getCurrentModuleIdx()
	if currentModuleIdx == -1 {
		return
	}

	for time.Now().Before(sh.Config.StartTime.Add(sh.Config.ModuleDuration * time.Duration(currentModuleIdx))) {
		time.Sleep(time.Until(sh.Config.StartTime.Add(sh.Config.ModuleDuration * time.Duration(currentModuleIdx))))
	}

	logger.Info.Printf("current module index: %d, starting in %f seconds\n", currentModuleIdx, time.Until(sh.Config.StartTime.Add(sh.Config.ModuleDuration*time.Duration(currentModuleIdx))).Seconds()*-1)

	if err := sh.launchModule(currentModuleIdx); err != nil {
		logger.Error.Printf("error launching module %02d: %v", currentModuleIdx, err)
		return
	}

	for {
		if time.Now().Before(sh.Config.StartTime.Add(sh.Config.ModuleDuration * time.Duration(currentModuleIdx+1))) {
			time.Sleep(time.Until(sh.Config.StartTime.Add(sh.Config.ModuleDuration * time.Duration(currentModuleIdx+1))))
		} else {
			currentModuleIdx += 1
			if err := sh.launchModule(currentModuleIdx); err != nil {
				logger.Error.Printf("error launching module %02d: %v\n", currentModuleIdx, err)
				return
			} else if currentModuleIdx == len(sh.Config.Modules) {
				logger.Info.Printf("last (%dth) module launch, returning from scheduler\n", len(sh.Config.Modules))
				return
			}
		}
	}
}

func (sh *Short) Launch() (err error) {
	logger.Info.Println("launching Short now!")

	go sh.schedule()

	return nil
}
