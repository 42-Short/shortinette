package webhook

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/short"
)

type WebhookTestMode struct {
	MonitoringFunction func()
	CurrentModule      Module.Module
}

// Initializes the webhook TestMode, which triggers submission grading
// as soon as activity is recorded on a user's main branch.
func NewWebhookTestMode(currentModule Module.Module) WebhookTestMode {
	wt := WebhookTestMode{MonitoringFunction: nil, CurrentModule: currentModule}
	wt.MonitoringFunction = func() {
		http.HandleFunc("/webhook", wt.handleWebhook)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			return
		}
	}
	return wt
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Commit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
}

var (
	mu sync.Mutex
)

func (wt *WebhookTestMode) handleWebhook(w http.ResponseWriter, r *http.Request) {
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
	if payload.Ref == "refs/heads/main" && payload.Pusher.Name != os.Getenv("GITHUB_ADMIN") {
		if strings.ToLower(payload.Commit.Message) == "grademe" {
			logger.Info.Println("push event identified as submission, proceeding to grade..")
			mu.Lock()
			defer mu.Unlock()

			go func() {
				if err := short.GradeModule(wt.CurrentModule, payload.Repository.Name); err != nil {
					logger.Error.Printf("error grading module: %v", err)
				}
			}()
		}
	}
}

func (m WebhookTestMode) Run() {
	m.MonitoringFunction()
}
