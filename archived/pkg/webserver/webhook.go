//go:build ignore

// webhook provides functionality to monitor GitHub webhook events and trigger
// grading of student submissions based on push events to the main branch.
package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/short"
	"github.com/gin-gonic/gin"
)

// Webhook represents the state and behavior for the webhook test mode, which
// triggers submission grading as soon as activity is recorded on a user's main branch.
type Webhook struct {
	MonitoringFunction func()                   // MonitoringFunction is the function that starts the webhook server.
	Modules            map[string]Module.Module // Modules is a map of module names to their corresponding Module structs.
	CurrentModule      string                   // CurrentModule is the name of the module currently being graded.
}

// NewWebhook initializes and returns a Webhook instance, which triggers
// submission grading as soon as activity is recorded on a user's main branch.
//
//   - modules: A map of module names to Module structs.
//   - endpoint: The endpoint the webhook is to be sending payloads to, with no trailing slash (e.g., '/webhook')
//   - port: The port the webhook is to be sending payloads to, without ':' (e.g., '8080')
//
// Returns a pointer to the initialized Webhook.
func NewWebhook(modules map[string]Module.Module) (webhook *Webhook) {
	return &Webhook{
		Modules: modules,
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

var (
	mu    sync.Mutex // mutex to prevent concurrent grading processes from overlapping.
	queue []string
	sem   = make(chan struct{}, 3)
)

func (wh *Webhook) HandleWebhook(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "invalid request method"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read request body: %v", err)})
		return
	}

	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse request body: %v", err)})
		return
	}

	if payload.Ref == "refs/heads/main" && payload.Pusher.Name != os.Getenv("GITHUB_ADMIN") {
		if strings.ToLower(payload.Commit.Message) == "grademe" {
			logger.Info.Printf("push event on %s identified as submission, grading..", payload.Repository.Name)

			mu.Lock()
			queue = append(queue, payload.Repository.Name)
			mu.Unlock()

			go wh.processQueue()
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "not a grading request, ignoring"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "grading request added to queue"})
}

func (wh *Webhook) processQueue() {
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
			if err := short.GradeModule(wh.Modules[wh.CurrentModule], repo, true); err != nil {
				logger.Error.Printf("error grading module for %s: %v", repo, err)
				return
			}
		}(repoName)
		queue = queue[1:]
	}
}
