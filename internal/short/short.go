package short

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/git"
	Module "github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
)

type HourlyTestMode struct {
	Delay             int
	FrequenzyDuration int
}

type MainBranchTestMode struct {
}

// TODO: find better name
type TestMode struct {
	Hourly     *HourlyTestMode
	MainBranch *MainBranchTestMode
}

func NewHourlyTestMode(hourly *HourlyTestMode) TestMode {
	return TestMode{
		Hourly:     hourly,
		MainBranch: nil,
	}
}
func NewMainBranchTestMode(mainBranch *MainBranchTestMode) TestMode {
	return TestMode{
		Hourly:     nil,
		MainBranch: mainBranch,
	}
}

// type SubjectSupplyMode enum {

// }

type Short struct {
	Name     string
	TestMode TestMode
	// Is it one per day or are users automatically assigne if they have the previous Subject at XX% --SubjectSupplyMode
	// modules [excercises]
}

func gradeModule(module Module.Module, config Config) error {
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

func endModule(module Module.Module, config Config) error {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "read"); err != nil {
			return err
		}
		if err := gradeModule(module, config); err != nil {
			return err
		}
	}
	return nil
}

func startModule(module Module.Module, config Config) error {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		if err := git.Create(repoId); err != nil {
			return err
		}
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "push"); err != nil {
			return err
		}
		if err := git.UploadFile(repoId, "subjects/R00.md", "README.md", fmt.Sprintf("Subject for module %s. Good Luck!", module.Name)); err != nil {
			return err
		}
	}
	return nil
}

func Run() {
	config, err := getConfig()
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return
	}
	if err := startModule(*R00.R00(), *config); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return
	}
	if err := endModule(*R00.R00(), *config); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return
	}
}
