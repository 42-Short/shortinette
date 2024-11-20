package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	_, err = data.SeedDB(db)
	if err != nil {
		logger.Error.Fatalf("failed to seed DB: %v", err)
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
	participant := data.Participant{
		IntraLogin:  "foo",
		GitHubLogin: "bar",
	}
	participantJson, err := json.Marshal(participant)
	require.NoError(t, err, "failed to marshal module")

	response := serveRequest(t, "POST", "/shortinette/v1/participants", strings.NewReader(string(participantJson)))
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, string(participantJson), response.Body.String())

}

func TestGetAll(t *testing.T) {
	response := serveRequest(t, "GET", "/shortinette/v1/participants", nil)

	var actualParticipants []data.Participant
	err := json.Unmarshal(response.Body.Bytes(), &actualParticipants)
	require.NoError(t, err, "failed to unmarshal module")

	participantDAO := data.NewDAO[data.Participant](api.DB)
	expectedParticipants, err := participantDAO.GetAll(context.Background())
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.Equal(t, expectedParticipants, actualParticipants)
}

func serveRequest(t *testing.T, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err, fmt.Sprintf("failed to make request: %s", url))

	req.Header.Set("Authorization", "Bearer "+apiToken)

	response := httptest.NewRecorder()
	api.Engine.ServeHTTP(response, req)
	return response
}
