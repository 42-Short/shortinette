package webhook

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	"github.com/42-Short/shortinette/pkg/short"
)

// Initializes the webhook TestMode, which triggers submission grading
// as soon as activity is recorded on a user's main branch.
func NewWebhookTestMode() WebhookTestMode {
	return WebhookTestMode{
		MonitoringFunction: func() {
			http.HandleFunc("/webhook", handleWebhook)
			if err := http.ListenAndServe(":8080", nil); err != nil {
				logger.Error.Printf("failed to start http server: %v", err)
			}
		},
	}
}

type WebhookTestMode struct {
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

var (
	mu sync.Mutex
)

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

	if payload.Ref == "refs/heads/main" && payload.Pusher.Name != os.Getenv("GITHUB_ADMIN") {
		logger.Info.Printf("received push event to main branch of %s by %s\n", payload.Repository.Name, payload.Pusher.Name)

		mu.Lock()
		defer mu.Unlock()

		go func() {
			if err := short.GradeModule(*R00.R00(), payload.Repository.Name); err != nil {
				logger.Error.Printf("error grading module: %v", err)
			}
		}()
	}
}

func (m WebhookTestMode) Run() {
	m.MonitoringFunction()
}
