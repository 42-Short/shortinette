package short

import (
	"fmt"
	"sort"
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

func updateRelease(repo db.Repository, newWaitingTime time.Duration, tracesPath string) error {
	releaseName := fmt.Sprintf("%d/100 - retry in %dm", repo.Score, int(newWaitingTime.Minutes()))

	if err := git.NewRelease(repo.ID, "Grade", releaseName, tracesPath, true); err != nil {
		return err
	}
	return nil
}

var tableRow = `<tr>
	<th>%s</th>
	<th>%s</th>
</tr>`

func getUpdatedReadme(repo db.Repository, results map[string]bool) (newReadme string, err error) {
	var keys []string
	for key := range results {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	oldContent, err := git.GetDecodedFile(repo.ID, "traces", "README.md")
	if err != nil {
		logger.Info.Printf("README.md not found in %s", repo.ID)
		oldContent = ""
	}

	var currentResult string
	newReadme = fmt.Sprintf("%s<h1 align=\"center\">ATTEMPT %d - SCORE %d/100</h1><div align=\"center\"><table>", oldContent, repo.Attempts, repo.Score)
	for _, key := range keys {
		if results[key] {
			currentResult = "OK"
		} else {
			currentResult = "KO"
		}
		newReadme = fmt.Sprintf("%s%s", newReadme, fmt.Sprintf(tableRow, key, currentResult))
	}
	newReadme = fmt.Sprintf("%s</table></div>", newReadme)
	return newReadme, nil
}

func uploadResults(repo db.Repository, tracesPath string, moduleName string, results map[string]bool) (err error) {
	commitMessage := fmt.Sprintf("Traces for module %s: %s", moduleName, tracesPath)
	if err := git.UploadFile(repo.ID, tracesPath, tracesPath, commitMessage, "traces"); err != nil {
		return err
	}

	updatedReadme, err := getUpdatedReadme(repo, results)
	if err != nil {
		return err
	}

	commitMessage = fmt.Sprintf("Results for module %s", moduleName)
	if err := git.UploadRaw(repo.ID, updatedReadme, "README.md", commitMessage, "traces"); err != nil {
		return err
	}
	return nil
}

func GradeModule(module Module.Module, repoID string) (err error) {
	repo, err := db.GetRepositoryData(module.Name, repoID)
	if err != nil {
		return err
	}
	repo.FirstAttempt = false

	if repo.WaitingTime > time.Since(repo.LastGradingTime) {
		logger.Info.Printf("repo %s attempted to grade too early", repo.ID)
		if err = updateRelease(repo, repo.WaitingTime-time.Since(repo.LastGradingTime), ""); err != nil {
			return err
		}
		return nil
	}

	results, tracesPath := module.Run(repoID, "studentcode")

	score, passed := module.GetScore(results)
	if passed {
		repo.WaitingTime = 15 * time.Minute
	} else {
		repo.WaitingTime = min(repo.WaitingTime+15*time.Minute, 60*time.Minute)
	}

	repo.Score = score

	if err = uploadResults(repo, tracesPath, module.Name, results); err != nil {
		return err
	}

	if err = updateRelease(repo, repo.WaitingTime, tracesPath); err != nil {
		return err
	}

	if err = db.UpdateRepository(module.Name, repo); err != nil {
		logger.Error.Printf("could not update %s: %v", repo.ID, err)
	}
	return nil
}

// Grades all participant's modules and upload traces.
func GradeAll(module Module.Module, config Config) error {
	for _, participant := range config.Participants {
		repoID := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := GradeModule(module, repoID); err != nil {
			return err
		}
	}
	return nil
}

// Grades all repos from a module and removes write access for all participants.
func EndModule(module Module.Module, config Config) {
	for _, participant := range config.Participants {
		repoID := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := git.AddCollaborator(repoID, participant.GithubUserName, "read"); err != nil {
			logger.Error.Printf("error adding collaborator: %v", err)
		}
		if err := GradeAll(module, config); err != nil {
			logger.Error.Printf("error grading module: %v", err)
		}
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
			repoID := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
			if err := git.Create(repoID); err != nil {
				logger.Error.Printf("error creating git repository: %v", err)
			}
			if err := git.AddCollaborator(repoID, participant.GithubUserName, "push"); err != nil {
				logger.Error.Printf("error adding collaborator: %v", err)
			}
			if err := git.UploadFile(repoID, "subjects/R00.md", "subject/README.md", fmt.Sprintf("Subject for module %s. Good Luck!", module.Name), ""); err != nil {
				logger.Error.Printf("error uploading file: %v", err)
			}
		}(participant)
	}
	wg.Wait()
}
