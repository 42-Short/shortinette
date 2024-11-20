package api

import (
	"context"
	"encoding/json"
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

	"github.com/stretchr/testify/assert"

	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/db"
	"github.com/42-Short/shortinette/logger"
	"github.com/joho/godotenv"
)

var api *API
var apiToken string

func setupSignalHandler(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
}

func shutdown(ctx context.Context, errCh chan error) {
	select {
	case err := <-errCh:
		if err != nil {
			logger.Error.Fatalf("failed to run server: %v", err)
		}
	case <-ctx.Done():
		err := api.Shutdown()
		if err != nil {
			logger.Error.Printf("failed to shutdown server: %v", err)
		}
		os.Exit(1)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandler(cancel)

	apiToken = os.Getenv("API_TOKEN")
	api = NewAPI(os.Getenv("SERVER_ADDR"), db)
	errCh := api.Run()
	go shutdown(ctx, errCh)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPost(t *testing.T) {
	participant := newDummyParticipant()
	w := httptest.NewRecorder()

	participantJson, err := json.Marshal(participant)
	if err != nil {
		t.Fatalf("failed to marshal module: %v", err)
	}
	req, _ := http.NewRequest("POST", "/shortinette/v1/participants", strings.NewReader(string(participantJson)))
	req.Header.Set("Authorization", "Bearer "+apiToken)

	api.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(participantJson), w.Body.String())

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
