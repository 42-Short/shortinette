package short

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/git"
	Module "github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/logger"
	ITestMode "github.com/42-Short/shortinette/internal/short/testmodes"
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

func endModule(module Module.Module, config Config) {
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

func startModule(module Module.Module, config Config) {
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

// c := cron.New(cron.WithSeconds())

// if _, err = c.AddFunc("0 * * * * ?", func() {
// 	module := R00.R00()
// 	logger.Info.Printf("starting module %s", module.Name)
// 	startModule(*module, *config)
// }); err != nil {
// 	logger.Error.Printf("failed scheduling start module task: %v", err)
// 	return
// }
// if _, err = c.AddFunc("59 * * * * ?", func() {
// 	module := R00.R00()
// 	logger.Info.Printf("ending module %s", module.Name)
// 	endModule(*module, *config)
// }); err != nil {
// 	logger.Error.Printf("failed scheduling end module task: %v", err)
// 	return
// }
// c.Start()
// select {}
