package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
	"github.com/joho/godotenv"
)

var api *API
var apiToken string

func shutdown(sigCh chan os.Signal, errCh chan error) {
	select {
	case err := <-errCh:
		if err != nil {
			logger.Error.Fatalf("failed to run server: %v", err)
		}
	case sig := <-sigCh:
		err := api.Shutdown()
		if err != nil {
			logger.Error.Printf("failed to shutdown server: %v", err)
		}
		logger.Error.Fatalf("caught signal: %v", sig)
	}
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Error.Fatalf("failed to load .env file: %v", err)
	}

	db, err := db.NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		logger.Error.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	err = db.Initialize("../db/schema.sql")
	if err != nil {
		logger.Error.Fatalf("failed to initialize db: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	api, err = NewAPI(db, gin.TestMode)
	if err != nil {
		logger.Error.Fatalf("failed to initialize API: %v", err)
	}
	errCh := api.Run()
	go shutdown(sigCh, errCh)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPost(t *testing.T) {
	participant := newDummyParticipant()
	participantJson, err := json.Marshal(participant)
	if err != nil {
		t.Fatalf("failed to marshal module: %v", err)
	}

	response := serveRequest(t, "POST", "/shortinette/v1/participants", strings.NewReader(string(participantJson)))
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, string(participantJson), response.Body.String())

}

// Todo: seed database instead of using POST
func TestGetAll(t *testing.T) {
	participant := newDummyParticipant()
	participantJson, err := json.Marshal(participant)
	if err != nil {
		t.Fatalf("failed to marshal module: %v", err)
	}

	response := serveRequest(t, "POST", "/shortinette/v1/participants", strings.NewReader(string(participantJson)))
	assert.Equal(t, http.StatusCreated, response.Code)

	response = serveRequest(t, "GET", "/shortinette/v1/participants", nil)
	assert.Equal(t, http.StatusOK, response.Code, response.Body)
}

func serveRequest(t *testing.T, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)

	response := httptest.NewRecorder()
	api.Engine.ServeHTTP(response, req)
	return response
}

func newDummyModule(moduleID int, intraLogin string) *data.Module {
	return &data.Module{
		Id:             moduleID,
		IntraLogin:     intraLogin,
		Attempts:       rand.Int(),
		Score:          rand.Int(),
		LastGraded:     time.Now(),
		WaitTime:       rand.Int(),
		GradingOngoing: rand.Intn(2) == 0,
	}
}

func newDummyParticipant() *data.Participant {
	intraLogin := strconv.Itoa(rand.Int())
	return &data.Participant{
		IntraLogin:  intraLogin,
		GitHubLogin: "dummy_git_" + intraLogin,
	}
}
