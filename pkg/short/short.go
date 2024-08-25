// Package short provides the core functionality for managing and grading coding modules.
// It handles initialization, grading, result uploading, and module lifecycle management.
package short

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/42-Short/shortinette/pkg/db"
	"github.com/42-Short/shortinette/pkg/git"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/requirements"
	ITestMode "github.com/42-Short/shortinette/pkg/short/testmodes"
)

// Short represents the main structure for managing a coding module, including the module's
// name, its exercises, and the test mode to use.
type Short struct {
	Name     string                   // Name is the display name of the Short.
	Modules  map[string]Module.Module // Modules is a map of module names to their corresponding Module structs.
	TestMode ITestMode.ITestMode      // TestMode determines how the submission testing will be triggered.
}

// shortInit initializes the logging and requirement validation for the Short application.
func shortInit() {
	logger.InitializeStandardLoggers("")
	if err := requirements.ValidateRequirements(); len(os.Args) != 2 && err != nil {
		logger.Error.Println(err.Error())
		return
	}
}

// NewShort returns a Short object, which serves as the wrapper for the entire Short
// configuration.
//
//   - name: the display name of your Short
//   - modules: a map of strings to Module.Module objects, used for quicker lookups during grading
//   - testMode: a ITestMode object, determining how the submission testing will be triggered
func NewShort(name string, modules map[string]Module.Module, testMode ITestMode.ITestMode) Short {
	shortInit()
	return Short{
		Name:     name,
		Modules:  modules,
		TestMode: testMode,
	}
}

// updateRelease updates the release information on the repository with the current grading
// results and the next grading attempt time.
//
//   - repo: the repository object containing grading details
//   - newWaitingTime: the duration to wait before the next grading attempt
//   - tracesPath: the path to the grading traces
//
// Returns an error if the release update fails.
func updateRelease(repo db.Repository, newWaitingTime time.Duration, tracesPath string) error {
	nextGradingAttemptTime := time.Now().Add(newWaitingTime).Format("15:04")
	releaseName := fmt.Sprintf("%d/100 - retry at %s", repo.Score, nextGradingAttemptTime)

	if err := git.NewRelease(repo.ID, "Grade", releaseName, tracesPath, true); err != nil {
		return err
	}
	return nil
}

// getUpdatedReadme generates the updated README content based on the latest grading
// results and appends it to the existing content.
//
//   - repo: the repository object containing grading details
//   - results: a map of exercise names to their pass/fail results
//
// Returns the new README content as a string and an error if the README update fails.
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
	var tableRow = `
<tr>
	<th>%s</th>
	<th>%s</th>
</tr>`
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

// uploadResults uploads the grading results and the updated README to the student's
// repository.
//
//   - repo: the repository object containing grading details
//   - tracesPath: the path to the grading traces
//   - moduleName: the name of the module being graded
//   - results: a map of exercise names to their pass/fail results
//
// Returns an error if the upload fails.
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

// checkPrematureGradingAttempt checks if a grading attempt is made before the waiting
// time has elapsed.
//
//   - repo: the repository object containing grading details
//
// Returns an error if the grading attempt is premature.
func checkPrematureGradingAttempt(repo db.Repository) (err error) {
	if os.Getenv("DEV_MODE") == "true" {
		return nil
	}
	if repo.WaitingTime > time.Since(repo.LastGradingTime) {
		if err = updateRelease(repo, repo.WaitingTime-time.Since(repo.LastGradingTime), ""); err != nil {
			return err
		}
		return fmt.Errorf("premature grading attempt")
	}
	return nil
}

// updateNewWaitingTime updates the waiting time for the next grading attempt based on
// the student's performance.
//
//   - repo: a pointer to the repository object containing grading details
//   - module: the module object containing the exercises
//   - results: a map of exercise names to their pass/fail results
func updateNewWaitingTime(repo *db.Repository, module Module.Module, results map[string]bool) {
	score, passed := module.GetScore(results)
	if passed {
		repo.WaitingTime = 15 * time.Minute
	} else {
		repo.WaitingTime = min(repo.WaitingTime+15*time.Minute, 60*time.Minute)
	}
	repo.Score = score
}

// GradeModule grades the exercises in a module for a specific student repository.
//
//   - module: the module object containing the exercises
//   - repoID: the ID of the student's repository
//
// Returns an error if the grading process fails.
func GradeModule(module Module.Module, repoID string) (err error) {
	repo, err := db.GetRepositoryData(module.Name, repoID)
	if err != nil {
		return err
	}
	repo.FirstAttempt = false

	if err = checkPrematureGradingAttempt(repo); err != nil {
		return err
	}

	results, tracesPath := module.Run(repoID)

	updateNewWaitingTime(&repo, module, results)

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

// GradeAll grades all participants' modules and uploads traces.
//
//   - module: the module object containing the exercises
//   - config: the configuration object containing participants' information
//
// Returns an error if the grading process fails for any participant.
func GradeAll(module Module.Module, config Config) error {
	for _, participant := range config.Participants {
		repoID := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := GradeModule(module, repoID); err != nil {
			return err
		}
	}
	return nil
}

// EndModule grades all repositories in a module and removes write access for all participants.
//
//   - module: the module object containing the exercises
//   - config: the configuration object containing participants' information
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

// StartModule creates a new repository for each participant, gives them write access,
// and uploads the module's subject to the repository.
//
//   - module: the module object containing the exercises
//   - config: the configuration object containing participants' information
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
			if err := git.Create(repoID, true, "traces"); err != nil {
				logger.Error.Printf("error creating git repository: %v", err)
			}
			if err := git.AddCollaborator(repoID, participant.GithubUserName, "push"); err != nil {
				logger.Error.Printf("error adding collaborator: %v", err)
			}
			if err := git.UploadFile(repoID, module.SubjectPath, "README.md", fmt.Sprintf("Subject for module %s. Good Luck!", module.Name), ""); err != nil {
				logger.Error.Printf("error uploading file: %v", err)
			}
		}(participant)
	}
	wg.Wait()
}

// dockerExecMode runs the grading process for a single exercise inside a Docker container.
//
//   - args: the command-line arguments passed to the application
//   - short: the Short object containing the module and test mode information
func dockerExecMode(short Short) {
	var config Module.GradingConfig
	err := json.Unmarshal([]byte(os.Args[1]), &config)
	if err != nil {
		os.Exit(1)
	}
	logger.InitializeStandardLoggers(config.ExerciseName)
	exercise, ok := short.Modules[config.ModuleName].Exercises[config.ExerciseName]
	if !ok {
		os.Exit(1)
	}
	exercise.CloneDirectory = config.TargetDirectory
	if err := logger.InitializeTraceLogger(config.TracesPath); err != nil {
		os.Exit(1)
	}
	result := exercise.Run()
	logger.File.Printf("[MOD%s][EX%s]: %s", config.ModuleName, config.ExerciseName, result.Output)
	if result.Passed {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Start begins the module lifecycle by starting the module and running the test mode.
//
//   - module: the name of the module to be started
func (short *Short) Start(module string) {
	config, err := GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	if len(os.Args) == 2 {
		dockerExecMode(*short)
	} else if len(os.Args) != 1 {
		logger.Error.Println("invalid number of arguments")
		return
	}
	StartModule(short.Modules[module], *config)
	short.TestMode.Run(module)
}
