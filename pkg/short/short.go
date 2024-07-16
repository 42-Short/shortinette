package short

import (
	"fmt"

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

func (h HourlyTestMode) Run() {
	// TODO
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

// Grades one participant's module and upload trace
func GradeModule(module Module.Module, repoId string) error {
	results, tracesPath := module.Run(repoId, "studentcode")
	commitMessage := fmt.Sprintf("Traces for module %s: %s", module.Name, tracesPath)
	if err := git.UploadFile(repoId, tracesPath, tracesPath, commitMessage); err != nil {
		return err
	}
	fmt.Println(results)
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
		if err := GradeAll(module, config); err != nil {
			logger.Error.Printf("error grading module: %v", err)
		}
	}
}

// Creates a new repo for each participant, gives them write access and
// uploads the module's subject on the repo.
func StartModule(module Module.Module, config Config) {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := git.Create(repoId); err != nil {
			logger.Error.Printf("error creating git repository: %v", err)
		}
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "push"); err != nil {
			logger.Error.Printf("error adding collaborator: %v", err)
		}
		if err := git.UploadFile(repoId, "subjects/R00.md", "README.md", fmt.Sprintf("Subject for module %s. Good Luck!", module.Name)); err != nil {
			logger.Error.Printf("error uploading file: %v", err)
		}
	}
}
