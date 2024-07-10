package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	"github.com/42-Short/shortinette/pkg/short"
)

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
	lastGradedTime time.Time
	mu             sync.Mutex
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
		fmt.Printf("Received push event to main branch of %s by %s\n", payload.Repository.Name, payload.Pusher.Name)
		config, err := short.GetConfig()
		if err != nil {
			http.Error(w, "failed to get module config", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if time.Since(lastGradedTime) < time.Minute {
			http.Error(w, "grading process is already running", http.StatusTooManyRequests)
			return
		}

		lastGradedTime = time.Now()

		go func() {
			if err := short.GradeModule(*R00.R00(), *config); err != nil {
				logger.Error.Printf("error grading module: %v", err)
			}
		}()
	}
}

func (m WebhookTestMode) Run() {
	m.MonitoringFunction()
}
