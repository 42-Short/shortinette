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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/dao"
	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
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
			logger.Error.Fatalf("failed to shutdown server: %v", err)
		}
		logger.Error.Fatalf("caught signal: %v", sig)
	}
}

func newDummyExercises() config.Exercise {
	ex, _ := config.NewExercise(10, []string{"foo.c", "bar.c"}, "foo")
	return *ex
}

func newDummyConfig() (*config.Config, error) {
	exercises := []config.Exercise{newDummyExercises(), newDummyExercises()}

	modules := []config.Module{
		{
			Exercises:    exercises,
			MinimumScore: 50,
		},
		{
			Exercises:    exercises,
			MinimumScore: 60,
		},
	}
	config := config.NewConfig(modules, time.Duration(24)*time.Hour, time.Now(), "", "")
	err := config.FetchEnvVariables()
	return config, err
}

func TestMain(m *testing.M) {
	os.Setenv("TZ", "UTC")
	_, err := time.LoadLocation("UTC")
	if err != nil {
		logger.Error.Fatalf("failed to load UTC location: %v", err)
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

	_, err = dao.SeedDB(db)
	if err != nil {
		logger.Error.Fatalf("failed to seed DB: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	config, err := newDummyConfig()
	if err != nil {
		logger.Error.Fatalf("failed to create dummy config: %v", err)
	}

	apiToken = config.ApiToken

	api = NewAPI(config, db, gin.TestMode)
	api.SetupRouter()

	errCh := make(chan error, 1)
	go func() {
		err = api.Run()
		errCh <- err
	}()
	go shutdown(sigCh, errCh)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestWebhookGrademe(t *testing.T) {
	payload := gitHubWebhookPayload{}
	payload.Ref = "refs/heads/main"
	payload.Repository.Name = "intra_login02"
	payload.Pusher.Name = "github_login"
	payload.Commit.Message = "grademe"

	payloadJson, err := json.Marshal(payload)
	require.NoError(t, err, "failed to marshal item")

	response := serveRequest(t, "POST", "http://localhost:8080/shortinette/v1/webhook/grademe", strings.NewReader(string(payloadJson)), apiToken)
	assert.Equal(t, http.StatusProcessing, response.Code, response.Body)
}

func TestGrademe(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	url := fmt.Sprintf("http://localhost:8080/shortinette/v1/modules/%d/%s/grademe", moduleID, intraLogin)
	response := serveRequest(t, "POST", url, nil, apiToken)
	assert.Equal(t, http.StatusProcessing, response.Code, response.Body)
}

func TestPostParticipant(t *testing.T) {
	testPost(t, dao.NewDummyParticipant(42), "/shortinette/v1/participants")
}

func TestPostModule(t *testing.T) {
	testPost(t, dao.NewDummyModule(42, "dummy_participant5"), "/shortinette/v1/modules")
}

func TestGetAllParticipants(t *testing.T) {
	testGetAll[dao.Participant](t, "/shortinette/v1/participants")
}

func TestGetAllModules(t *testing.T) {
	testGetAll[dao.Module](t, "/shortinette/v1/modules")
}

func TestGetModule(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	url := fmt.Sprintf("/shortinette/v1/modules/%d/%s", moduleID, intraLogin)
	testGet[dao.Module](t, url, moduleID, intraLogin)
}

func TestGetParticipant(t *testing.T) {
	const intraLogin = "dummy_participant5"
	url := fmt.Sprintf("/shortinette/v1/participants/%s", intraLogin)
	testGet[dao.Participant](t, url, intraLogin)
}

func TestPutParticipant(t *testing.T) {
	const intraLogin = "dummy_participant5"
	testPut[dao.Participant](t, "/shortinette/v1/participants", intraLogin)
}

func TestPutModule(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	testPut[dao.Module](t, "/shortinette/v1/modules", moduleID, intraLogin)
}

func TestDeleteParticipant(t *testing.T) {
	const intraLogin = "dummy_participant5"
	url := fmt.Sprintf("/shortinette/v1/participants/%s", intraLogin)
	testDelete[dao.Participant](t, url, intraLogin)
}

func TestDeleteModule(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	url := fmt.Sprintf("/shortinette/v1/modules/%d/%s", moduleID, intraLogin)
	testDelete[dao.Module](t, url, intraLogin)
}

func TestUnauthorized(t *testing.T) {
	response := serveRequest(t, "GET", "/shortinette/v1/participants", nil, "foo")

	assert.Equal(t, http.StatusUnauthorized, response.Code, response.Body)
}

func testPost(t *testing.T, item any, url string) {
	t.Helper()

	itemJson, err := json.Marshal(item)
	require.NoError(t, err, "failed to marshal item")

	response := serveRequest(t, "POST", url, strings.NewReader(string(itemJson)), apiToken)
	assert.Equal(t, http.StatusCreated, response.Code, response.Body)
	assert.Equal(t, string(itemJson), response.Body.String())
}

func testPut[T any](t *testing.T, url string, args ...any) {
	t.Helper()

	dao := dao.NewDAO[T](api.DB)
	item, err := dao.Get(context.Background(), args...)
	require.NoError(t, err)

	itemJson, err := json.Marshal(item)
	require.NoError(t, err, "failed to marshal item")

	response := serveRequest(t, "PUT", url, strings.NewReader(string(itemJson)), apiToken)
	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.NotEqual(t, string(itemJson), response.Body.String())
}

func testGetAll[T any](t *testing.T, url string) {
	t.Helper()

	response := serveRequest(t, "GET", url, nil, apiToken)

	var actualItems []T
	err := json.Unmarshal(response.Body.Bytes(), &actualItems)
	require.NoError(t, err, "failed to unmarshal item")

	dao := dao.NewDAO[T](api.DB)
	expectedItems, err := dao.GetAll(context.Background())
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.Equal(t, expectedItems, actualItems)
}

func testGet[T any](t *testing.T, url string, args ...any) {
	t.Helper()

	response := serveRequest(t, "GET", url, nil, apiToken)

	var actualItem T
	err := json.Unmarshal(response.Body.Bytes(), &actualItem)
	require.NoError(t, err, "failed to unmarshal module")

	dao := dao.NewDAO[T](api.DB)
	expectedItem, err := dao.Get(context.Background(), args...)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.Equal(t, *expectedItem, actualItem)
}

func testDelete[T any](t *testing.T, url string, args ...any) {
	t.Helper()

	response := serveRequest(t, "DELETE", url, nil, apiToken)
	assert.Equal(t, http.StatusOK, response.Code, response.Body)

	dao := dao.NewDAO[T](api.DB)
	_, err := dao.Get(context.Background(), args...)
	require.Error(t, err)
}

func serveRequest(t *testing.T, method string, url string, body io.Reader, accessToken string) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err, fmt.Sprintf("failed to make request: %s", url))

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	response := httptest.NewRecorder()
	api.Engine.ServeHTTP(response, req)
	return response
}
