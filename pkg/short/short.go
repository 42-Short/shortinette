package short

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/pkg/git"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	ITestMode "github.com/42-Short/shortinette/pkg/short/testmodes"
)

type HourlyTestMode struct {
	Delay              int
	FrequencyDuration  int
	MonitoringFunction func()
}

type Repository struct {
	FirstAttempt    bool
	LastGradingTime time.Time
	WaitingTime     time.Duration
}

type Short struct {
	Name     string
	Modules  map[string]Module.Module
	TestMode ITestMode.ITestMode
}

// Returns a Short object, the wrapper for the whole Short configuration.
//
//   - name: the display name of your Short
//   - testMode: a ITestMode object, determining how the submission testing will
//     be triggered
func NewShort(name string, modules map[string]Module.Module, testMode ITestMode.ITestMode) Short {
	return Short{
		Name:     name,
		Modules:  modules,
		TestMode: testMode,
	}
}

func getScore(results map[string]bool, module Module.Module) int {
	var keys []string
	for key := range results {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	score := 0

	for _, key := range keys {
		if !results[key] {
			break
		} else {
			score += module.Exercises[key].Score
		}
	}
	return score
}

func uploadScore(module Module.Module, repoId string, results map[string]bool, newWaitingTime time.Duration) error {
	score := getScore(results, module)
	releaseName := fmt.Sprintf("%d/100", score)

	if err := git.NewRelease(repoId, "Grade", releaseName, true, newWaitingTime); err != nil {
		return err
	}
	return nil
}

func GradeModule(module Module.Module, repoId string, repositories map[string]Repository) error {
	fmt.Println(repositories)
	repo := repositories[repoId]
	if repo.LastGradingTime.IsZero() {
		repo.LastGradingTime = time.Now()
		repo.FirstAttempt = false
	}

	oldScore := git.GetLatestScore(repoId)

	if repo.WaitingTime > time.Since(repo.LastGradingTime) {
		logger.Info.Printf("repo '%s' attempted grading too early", repoId)
		scoreString := fmt.Sprintf("%d/100", oldScore)
		waitingTime := time.Duration(repo.WaitingTime - time.Since(repo.LastGradingTime))
		if err := git.NewRelease(repoId, "Grade", scoreString, true, waitingTime); err != nil {
			return err
		}
	}

	results, tracesPath := module.Run(repoId, "studentcode")

	if getScore(results, module) > module.MinimumGrade {
		repo.WaitingTime = 15 * time.Minute
	} else {
		repo.WaitingTime = min(repo.WaitingTime+15*time.Minute, 60*time.Minute)
	}

	commitMessage := fmt.Sprintf("Traces for module %s: %s", module.Name, tracesPath)

	if err := git.UploadFile(repoId, tracesPath, tracesPath, commitMessage, "traces"); err != nil {
		return err
	}

	if err := uploadScore(module, repoId, results, repo.WaitingTime); err != nil {
		return err
	}

	repositories[repoId] = repo
	return nil
}

// Grades all participant's modules and upload traces.
func GradeAll(module Module.Module, config Config, repositories map[string]Repository) error {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := GradeModule(module, repoId, repositories); err != nil {
			return err
		}
	}
	return nil
}

// Grades all repos from a module and removes write access for all participants.
func EndModule(module Module.Module, config Config, repositories map[string]Repository) {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "read"); err != nil {
			logger.Error.Printf("error adding collaborator: %v", err)
		}
		if err := GradeAll(module, config, repositories); err != nil {
			logger.Error.Printf("error grading module: %v", err)
		}
	}
}

func newRepository() Repository {
	return Repository{
		FirstAttempt:    true,
		LastGradingTime: time.Time{},
		WaitingTime:     time.Duration(0),
	}
}

func GetRepositories(config Config, moduleName string) map[string]Repository {
	repositories := make(map[string]Repository)

	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, moduleName)
		repositories[repoId] = newRepository()
	}
	return repositories
}

// StartModule creates a new repo for each participant, gives them write access, and uploads the module's subject on the repo.
func StartModule(module Module.Module, config Config) {
	var wg sync.WaitGroup

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
