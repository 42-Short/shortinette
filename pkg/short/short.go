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
	TestMode ITestMode.ITestMode
}

func NewShort(name string, testMode ITestMode.ITestMode) Short {
	return Short{
		Name:     name,
		TestMode: testMode,
	}
}

func GradeModule(module Module.Module, config Config) error {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		result, tracesPath := module.Run(repoId, "studentcode")
		if err := git.UploadFile(repoId, tracesPath, tracesPath, fmt.Sprintf("Traces for module %s: %s", module.Name, tracesPath)); err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}

func EndModule(module Module.Module, config Config) {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		// INFO: Giving read access to a user will remove their push rights
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "read"); err != nil {
			logger.Error.Printf("error adding collaborator: %v", err)
		}
		if err := GradeModule(module, config); err != nil {
			logger.Error.Printf("error grading module: %v", err)
		}
	}
}

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
