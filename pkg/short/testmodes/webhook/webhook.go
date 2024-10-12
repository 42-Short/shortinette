// webhook provides functionality to monitor GitHub webhook events and trigger
// grading of student submissions based on push events to the main branch.
package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/short"
)

// WebhookTestMode represents the state and behavior for the webhook test mode, which
// triggers submission grading as soon as activity is recorded on a user's main branch.
type WebhookTestMode struct {
	MonitoringFunction func()                   // MonitoringFunction is the function that starts the webhook server.
	Modules            map[string]Module.Module // Modules is a map of module names to their corresponding Module structs.
	CurrentModule      string                   // CurrentModule is the name of the module currently being graded.
	server             *http.Server
	mu                 sync.Mutex
	endpoint           string
	port               string
}

// NewWebhookTestMode initializes and returns a WebhookTestMode instance, which triggers
// submission grading as soon as activity is recorded on a user's main branch.
//
//   - modules: A map of module names to Module structs.
//   - endpoint: The endpoint the webhook is to be sending payloads to, with no trailing slash (e.g., '/webhook')
//   - port: The port the webhook is to be sending payloads to, without ':' (e.g., '8080')
//
// Returns a pointer to the initialized WebhookTestMode.
func NewWebhookTestMode(modules map[string]Module.Module, endpoint string, port string) *WebhookTestMode {
	return &WebhookTestMode{
		Modules:  modules,
		endpoint: endpoint,
		port:     port,
	}
}

// GitHubWebhookPayload represents the structure of the JSON payload sent by GitHub when
// a push event occurs.
type GitHubWebhookPayload struct {
	Ref        string `json:"ref"` // Ref is the git reference (branch) that was pushed to.
	Repository struct {
		Name string `json:"name"` // Name is the name of the repository where the push occurred.
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"` // Name is the name of the user who pushed the commit.
	} `json:"pusher"`
	Commit struct {
		Message string `json:"message"` // Message is the commit message of the push.
	} `json:"head_commit"`
}

func (wt *WebhookTestMode) startServer() {
	wt.mu.Lock()
	defer wt.mu.Unlock()

	wt.stopServer()

	mux := http.NewServeMux()
	mux.HandleFunc(wt.endpoint, wt.handleWebhook)

	wt.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", wt.port),
		Handler: mux,
	}

	go func() {
		if err := wt.server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error.Printf("HTTP server ListenAndServe: %v\n", err)
		}
	}()

	logger.Info.Printf("Webhook server started for module %s\n", wt.CurrentModule)
}

func (wt *WebhookTestMode) stopServer() {
	if wt.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := wt.server.Shutdown(ctx); err != nil {
			logger.Error.Printf("Server forced to shut down: %v\n", err)
		}
		wt.server = nil
		logger.Info.Println("Webhook server stopped")
	}
}

var (
	mu    sync.Mutex // mutex to prevent concurrent grading processes from overlapping.
	queue []string
	sem   = make(chan struct{}, 3)
)

// handleWebhook processes incoming webhook events and triggers grading if the event
// corresponds to a push to the main branch with the commit message "grademe".
//
//   - w: The http.ResponseWriter used to send the response.
//   - r: The http.Request representing the incoming webhook event.
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
			logger.Info.Printf("push event on %s identified as submission, grading..", payload.Repository.Name)

			mu.Lock()
			queue = append(queue, payload.Repository.Name)
			mu.Unlock()

			go wt.processQueue()
		}
	}
}

func (wt *WebhookTestMode) processQueue() {
	for {
		mu.Lock()
		if len(queue) == 0 {
			mu.Unlock()
			return
		}
		repoName := queue[0]
		mu.Unlock()

		sem <- struct{}{}
		go func(repo string) {
			defer func() { <-sem }()
			if err := short.GradeModule(wt.Modules[wt.CurrentModule], repo, true); err != nil {
				logger.Error.Printf("error grading module for %s: %v", repo, err)
				return
			}
		}(repoName)
		queue = queue[1:]
	}
}

// Run starts the webhook server and sets the current module to be graded.
//
//   - currentModule: The name of the module that is currently being graded.
func (wt *WebhookTestMode) Run(currentModule string) {
	wt.mu.Lock()
	wt.CurrentModule = currentModule
	wt.mu.Unlock()
	wt.startServer()
}
