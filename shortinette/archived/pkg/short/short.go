//go:build ignore

// short provides the core functionality for managing and grading coding modules.
// It handles initialization, grading, result uploading, and module lifecycle management.
package short

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/42-Short/shortinette/pkg/db"
	"github.com/42-Short/shortinette/pkg/git"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/requirements"
)

// Short represents the main structure for managing a coding module, including the module's
// name, its exercises, and the test mode to use.
type Short struct {
	Name    string                   // Name is the display name of the Short.
	Modules map[string]Module.Module // Modules is a map of module names to their corresponding Module structs.
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
func NewShort(name string, modules map[string]Module.Module) Short {
	shortInit()
	return Short{
		Name:    name,
		Modules: modules,
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
		return fmt.Errorf("adding release to %s: %v", repo.ID, err)
	}
	return nil
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
func uploadResults(repo db.Repository, tracesPath string, moduleName string) (err error) {
	commitMessage := fmt.Sprintf("Traces for module %s: %s", moduleName, tracesPath)
	if err := git.UploadFile(repo.ID, tracesPath, tracesPath, commitMessage, "traces"); err != nil {
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
	repo.Score = max(score, repo.Score)
}

// Sorts the trace content before uploading (containers are writing into the file asynchronously,
// leading to mixups in the output)
func sortTraceContent(tracesPath string) (err error) {
	contentAsBytes, err := os.ReadFile(tracesPath)
	if err != nil {
		return err
	}
	outputByExercise := make(map[int][]string)
	contentAsSlice := strings.Split(string(contentAsBytes), "\n")
	pattern := regexp.MustCompile(`\[MOD\d+\]\[EX(\d+)\]`)
	var exerciseNumber int
	for _, line := range contentAsSlice {
		if match := pattern.FindStringSubmatch(line); len(match) > 1 {
			exerciseNumber, _ = strconv.Atoi(match[1])
		}
		if line != "" {
			outputByExercise[exerciseNumber] = append(outputByExercise[exerciseNumber], line)
		}
	}

	var keys []int
	for k := range outputByExercise {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var sortedOutput strings.Builder
	for _, k := range keys {
		for _, line := range outputByExercise[k] {
			sortedOutput.WriteString(line + "\n")
		}
	}

	err = os.WriteFile(tracesPath, []byte(sortedOutput.String()), 0o644)
	if err != nil {
		return err
	}

	return nil
}

// GradeModule grades the exercises in a module for a specific student repository.
//
//   - module: the module object containing the exercises
//   - repoID: the ID of the student's repository
//
// Returns an error if the grading process fails.
func GradeModule(module Module.Module, repoID string, updateDatabase bool) (err error) {
	repo, err := db.GetRepositoryData(module.Name, repoID)
	if err != nil {
		return fmt.Errorf("could not get repository data: %v", err)
	}
	repo.FirstAttempt = false

	if updateDatabase {
		defer func() {
			if err = db.UpdateRemoteDatabase(); err != nil {
				logger.Error.Printf("failure to update remote database: %v", err)
			}
		}()
	}

	if err = checkPrematureGradingAttempt(repo); err != nil {
		return err
	}

	results, tracesPath, err := module.Run(repoID)
	if err != nil {
		return err
	}

	updateNewWaitingTime(&repo, module, results)

	if err := sortTraceContent(tracesPath); err != nil {
		return fmt.Errorf("sorting trace content: %v", err)
	}

	if err = uploadResults(repo, tracesPath, module.Name); err != nil {
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

const maxConcurrentGET = 5
const maxConcurrentPOST = 1

// EndModule grades all repositories in a module and removes write access for all participants.
//
//   - module: the module object containing the exercises
//   - config: the configuration object containing participants' information
func EndModule(module Module.Module, config Config) (err error) {
	defer func() {
		if err = db.UpdateRemoteDatabase(); err != nil {
			logger.Error.Printf("updating remote database: %v", err)
		}
	}()
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentGET)
	errChan := make(chan error, len(config.Participants))

	for _, participant := range config.Participants {
		wg.Add(1)
		sem <- struct{}{}
		repoID := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)

		go func(repoId string, module Module.Module, githubUser string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := git.AddCollaborator(repoID, participant.GithubUserName, "read"); err != nil {
				errChan <- fmt.Errorf("error adding collaborator: %v", err)
				return
			}
			if err := GradeModule(module, repoID, false); err != nil {
				errChan <- fmt.Errorf("error grading module: %v", err)
				return
			}
		}(repoID, module, participant.GithubUserName)
	}
	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to grade %d repositories: %v", len(errors), errors)
	}
	return nil
}

// Asynchronously creates a repo for each user in the config specified by CONFIG_PATH.
//
//   - config: Config struct filled with the participant's data
//   - module: Module.Module struct filled with the module's metadata
func initializeRepos(config Config, module Module.Module) (err error) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentPOST)
	errChan := make(chan error, len(config.Participants))

	for _, participant := range config.Participants {
		wg.Add(1)
		sem <- struct{}{}

		go func(participant Participant) {
			defer wg.Done()
			defer func() { <-sem }()

			repoID := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
			if err := git.Create(repoID, true, "traces"); err != nil {
				errChan <- fmt.Errorf("error creating git repository: %v", err)
				return
			}
			if err := git.AddCollaborator(repoID, participant.GithubUserName, "push"); err != nil {
				errChan <- fmt.Errorf("error adding collaborator: %v", err)
				return
			}
			if err := git.UploadFile(repoID, module.SubjectPath, "README.md", fmt.Sprintf("Subject for module %s. Good Luck!", module.Name), ""); err != nil {
				errChan <- fmt.Errorf("error uploading file: %v", err)
				return
			}
		}(participant)
	}
	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to initialize %d repositories: %v", len(errors), errors)
	}
	return nil
}

// StartModule creates a new repository for each participant, gives them write access,
// and uploads the module's subject to the repository.
//
//   - module: the module object containing the exercises
//   - config: the configuration object containing participants' information
func StartModule(module Module.Module, config Config) (err error) {
	defer func() {
		if err = db.UpdateRemoteDatabase(); err != nil {
			logger.Error.Printf("updating remote database: %v", err)
		}
	}()

	if err = db.CreateTable(fmt.Sprintf("repositories_%s", module.Name)); err != nil {
		return fmt.Errorf("table creation: %v", err)
	}

	participants := [][]string{}
	for _, participant := range config.Participants {
		participants = append(participants, []string{participant.GithubUserName, participant.IntraLogin})
	}

	if err = db.InitModuleTable(participants, module.Name); err != nil {
		return fmt.Errorf("table initializtion: %v", err)
	}

	if err := initializeRepos(config, module); err != nil {
		return fmt.Errorf("repo initialization: %v", err)
	}

	return nil
}

// dockerExecMode runs the grading process for a single exercise inside a Docker container.
// Here, we unmarshal a JSON passed by the main program through command line arguments, containing
// all the state information we need for grading the exercise.
//
//   - short: the Short object containing the module and test mode information
func dockerExecMode(short Short) {
	var config Module.GradingConfig
	err := json.Unmarshal([]byte(os.Args[1]), &config)
	if err != nil {
		logger.Error.Printf("%s is not a valid Module.GradingConfig struct and cannot be unmarshalled by shortinette.", os.Args[1])
		os.Exit(1)
	}
	logger.InitializeStandardLoggers(config.ExerciseName)
	exercise, ok := short.Modules[config.ModuleName].Exercises[config.ExerciseName]
	if !ok {
		logger.Error.Printf("module %s, exercise %s not found - fix your GradingConfig struct or add it to your Short.", config.ModuleName, config.ExerciseName)
		os.Exit(1)
	}
	exercise.CloneDirectory = config.CloneDirectory
	if err := logger.InitializeTraceLogger(config.TracesPath); err != nil {
		logger.Error.Printf("logger initialization: %v", err)
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
func (short *Short) StartModule(module string) (err error) {
	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	if err = StartModule(short.Modules[module], *config); err != nil {
		return fmt.Errorf("error starting module: %v", err)
	}

	return nil
}

func (short *Short) Start() {
	if len(os.Args) == 2 {
		dockerExecMode(*short)
	}
}
