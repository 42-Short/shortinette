package short

import (
	"fmt"
	"sync"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/pkg/db"
	"github.com/42-Short/shortinette/pkg/git"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	ITestMode "github.com/42-Short/shortinette/pkg/short/testmodes"
)

type HourlyTestMode struct {
	Delay              int
	FrequencyDuration  int
	MonitoringFunction func()
}

type Short struct {
	Name     string
	Modules  map[string]Module.Module
	TestMode ITestMode.ITestMode
}

// Returns a Short object, the wrapper for the whole Short configuration.
//
//   - name: the display name of your Short
//   - modules: a map of strings to Module.Module objects, used for quicker lookups during grading
//   - testMode: a ITestMode object, determining how the submission testing will
//     be triggered
func NewShort(name string, modules map[string]Module.Module, testMode ITestMode.ITestMode) Short {
	return Short{
		Name:     name,
		Modules:  modules,
		TestMode: testMode,
	}
}

func updateRelease(repoId string, newWaitingTime time.Duration, score int) error {
	releaseName := fmt.Sprintf("%d/100 - retry in %dm", score, int(newWaitingTime.Minutes()))

	if err := git.NewRelease(repoId, "Grade", releaseName, true); err != nil {
		return err
	}
	return nil
}

func GradeModule(module Module.Module, repoId string) (err error) {
	repo, err := db.GetRepositoryData(module.Name, repoId)
	if err != nil {
		return err
	}
	repo.FirstAttempt = false

	oldScore := git.GetLatestScore(repoId)

	if repo.WaitingTime > time.Since(repo.LastGradingTime) {
		logger.Info.Printf("repo %s attempted to grade too early", repo.ID)
		if err = updateRelease(repo.ID, repo.WaitingTime-time.Since(repo.LastGradingTime) - 2 * time.Hour, oldScore); err != nil {
			return err
		}
		return nil
	}

	results, tracesPath := module.Run(repoId, "studentcode")

	score, passed := module.GetScore(results)
	if passed {
		repo.WaitingTime = 15 * time.Minute
	} else {
		repo.WaitingTime = min(repo.WaitingTime+15*time.Minute, 60*time.Minute)
	}

	commitMessage := fmt.Sprintf("Traces for module %s: %s", module.Name, tracesPath)

	if err := git.UploadFile(repoId, tracesPath, tracesPath, commitMessage, "traces"); err != nil {
		return err
	}

	if err := updateRelease(repoId, repo.WaitingTime, score); err != nil {
		return err
	}
	repo.Score = score
	if err = db.UpdateRepository(module.Name, repo); err != nil {
		logger.Error.Printf("could not update %s: %v", repo.ID, err)
	}
	return nil
}

// Grades all participant's modules and upload traces.
func GradeAll(module Module.Module, config Config) error {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := GradeModule(module, repoId); err != nil {
			return err
		}
	}
	return nil
}

// Grades all repos from a module and removes write access for all participants.
func EndModule(module Module.Module, config Config) {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "read"); err != nil {
			logger.Error.Printf("error adding collaborator: %v", err)
		}
		// if err := GradeAll(module, config); err != nil {
		// 	logger.Error.Printf("error grading module: %v", err)
		// }
	}
}

// StartModule creates a new repo for each participant, gives them write access, and uploads the module's subject on the repo.
func StartModule(module Module.Module, config Config) {
	var wg sync.WaitGroup

	created, err := db.CreateTable(fmt.Sprintf("repositories_%s", module.Name))
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	if created {
		var participants = [][]string{}
		for _, participant := range config.Participants {
			participants = append(participants, []string{participant.GithubUserName, participant.IntraLogin})
		}
		if err := db.InitModuleTable(participants, module.Name); err != nil {
			logger.Error.Println(err.Error())
			return
		}
	}

	for _, participant := range config.Participants {
		wg.Add(1)
		go func(participant Participant) {
			defer wg.Done()
			repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
			if err := git.Create(repoId); err != nil {
				logger.Error.Printf("error creating git repository: %v", err)
			}
			if err := git.AddCollaborator(repoId, participant.GithubUserName, "push"); err != nil {
				logger.Error.Printf("error adding collaborator: %v", err)
			}
			if err := git.UploadFile(repoId, "subjects/R00.md", "subject/README.md", fmt.Sprintf("Subject for module %s. Good Luck!", module.Name), ""); err != nil {
				logger.Error.Printf("error uploading file: %v", err)
			}
		}(participant)
	}
	wg.Wait()
}
