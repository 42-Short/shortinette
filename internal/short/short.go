package short

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/42-Short/shortinette/internal/git"
	Module "github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
)

type ITestMode interface {
	Run()
}

type HourlyTestMode struct {
	Delay              int
	FrequencyDuration  int
	MonitoringFunction func()
}

func (h HourlyTestMode) Run() {
	// TODO
}

type MainBranchTestMode struct {
	MonitoringFunction func()
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "failed to parse request body", http.StatusInternalServerError)
		return
	}

	if payload.Ref == "refs/heads/main" {
		fmt.Printf("Received push event to main branch of %s by %s\n", payload.Repository.Name, payload.Pusher.Name)
		config, err := getConfig()
		if err != nil {
			http.Error(w, "failed to get module config", http.StatusInternalServerError)
			return
		}
		go gradeModule(*R00.R00(), *config)
	}
}

func (m MainBranchTestMode) Run() {
	// TODO
}

type Short struct {
	Name     string
	TestMode ITestMode
}

func NewShort(name string, testMode ITestMode) Short {
	return Short{
		Name:     name,
		TestMode: testMode,
	}
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

func endModule(module Module.Module, config Config) {
	for _, participant := range config.Participants {
		repoId := fmt.Sprintf("%s-%s", participant.IntraLogin, module.Name)
		// INFO: Giving read access to a user will remove their push rights
		if err := git.AddCollaborator(repoId, participant.GithubUserName, "read"); err != nil {
			logger.Error.Printf("error adding collaborator: %v", err)
		}
		if err := gradeModule(module, config); err != nil {
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

func Run() {
	http.HandleFunc("/webhook", handleWebhook)
	http.ListenAndServe(":8080", nil)
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
