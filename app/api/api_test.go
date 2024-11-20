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

	"github.com/42-Short/shortinette/data"
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

func TestMain(m *testing.M) {
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

	api = NewAPI(db, gin.TestMode, time.Minute)
	errCh := api.Run()
	go shutdown(sigCh, errCh)

	apiToken = api.accessToken

	exitCode := m.Run()
	os.Exit(exitCode)
}

//TODO: new seeded db for every test

func TestPostParticipant(t *testing.T) {
	testPost(t, data.NewDummyParticipant(42), "/shortinette/v1/participants")
}

func TestPostModule(t *testing.T) {
	testPost(t, data.NewDummyModule(42, "dummy_participant5"), "/shortinette/v1/modules")
}

func TestGetAllParticipants(t *testing.T) {
	testGetAll[data.Participant](t, "/shortinette/v1/participants")
}

func TestGetAllModules(t *testing.T) {
	testGetAll[data.Module](t, "/shortinette/v1/modules")
}

func TestGetModule(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	url := fmt.Sprintf("/shortinette/v1/modules/%d/%s", moduleID, intraLogin)
	testGet[data.Module](t, url, moduleID, intraLogin)
}

func TestGetParticipant(t *testing.T) {
	const intraLogin = "dummy_participant5"
	url := fmt.Sprintf("/shortinette/v1/participants/%s", intraLogin)
	testGet[data.Participant](t, url, intraLogin)
}

func TestPutParticipant(t *testing.T) {
	const intraLogin = "dummy_participant5"
	testPut[data.Participant](t, "/shortinette/v1/participants", intraLogin)
}

func TestPutModule(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	testPut[data.Module](t, "/shortinette/v1/modules", moduleID, intraLogin)
}

func TestDeleteParticipant(t *testing.T) {
	const intraLogin = "dummy_participant5"
	url := fmt.Sprintf("/shortinette/v1/participants/%s", intraLogin)
	testDelete[data.Participant](t, url, intraLogin)
}

func TestDeleteModule(t *testing.T) {
	const (
		intraLogin = "dummy_participant5"
		moduleID   = 0
	)
	url := fmt.Sprintf("/shortinette/v1/modules/%d/%s", moduleID, intraLogin)
	testDelete[data.Module](t, url, intraLogin)
}

func testPost(t *testing.T, item any, url string) {
	t.Helper()

	itemJson, err := json.Marshal(item)
	require.NoError(t, err, "failed to marshal item")

	response := serveRequest(t, "POST", url, strings.NewReader(string(itemJson)))
	assert.Equal(t, http.StatusCreated, response.Code, response.Body)
	assert.Equal(t, string(itemJson), response.Body.String())
}

func testPut[T any](t *testing.T, url string, args ...any) {
	t.Helper()

	dao := data.NewDAO[T](api.DB)
	item, err := dao.Get(context.Background(), args...)
	require.NoError(t, err)

	itemJson, err := json.Marshal(item)
	require.NoError(t, err, "failed to marshal item")

	response := serveRequest(t, "PUT", url, strings.NewReader(string(itemJson)))
	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.NotEqual(t, string(itemJson), response.Body.String())
}

func testGetAll[T any](t *testing.T, url string) {
	t.Helper()

	response := serveRequest(t, "GET", url, nil)

	var actualItems []T
	err := json.Unmarshal(response.Body.Bytes(), &actualItems)
	require.NoError(t, err, "failed to unmarshal item")

	dao := data.NewDAO[T](api.DB)
	expectedItems, err := dao.GetAll(context.Background())
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.Equal(t, expectedItems, actualItems)
}

func testGet[T any](t *testing.T, url string, args ...any) {
	t.Helper()

	response := serveRequest(t, "GET", url, nil)

	var actualItem T
	err := json.Unmarshal(response.Body.Bytes(), &actualItem)
	require.NoError(t, err, "failed to unmarshal module")

	dao := data.NewDAO[T](api.DB)
	expectedItem, err := dao.Get(context.Background(), args...)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Code, response.Body)
	assert.Equal(t, *expectedItem, actualItem)
}

func testDelete[T any](t *testing.T, url string, args ...any) {
	t.Helper()

	response := serveRequest(t, "DELETE", url, nil)
	assert.Equal(t, http.StatusOK, response.Code, response.Body)

	dao := data.NewDAO[T](api.DB)
	_, err := dao.Get(context.Background(), args...)
	require.Error(t, err)
}

func serveRequest(t *testing.T, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err, fmt.Sprintf("failed to make request: %s", url))

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	response := httptest.NewRecorder()
	api.Engine.ServeHTTP(response, req)
	return response
}
