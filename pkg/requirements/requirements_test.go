package requirements

import (
	"os"
	"testing"
)

func TestMissingGithubAdmin(t *testing.T) {
	os.Unsetenv("GITHUB_ADMIN")
	if err := requireEnv(); err == nil {
		t.Fatalf("missing required environment variable not caught")
	}
}

func TestMissingGithubToken(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")
	if err := requireEnv(); err == nil {
		t.Fatalf("missing required environment variable not caught")
	}
}

func TestMissingGithubOrga(t *testing.T) {
	os.Unsetenv("GITHUB_ORGANISATION")
	if err := requireEnv(); err == nil {
		t.Fatalf("missing required environment variable not caught")
	}
}

func TestMissingGithubEmail(t *testing.T) {
	os.Unsetenv("GITHUB_EMAIL")
	if err := requireEnv(); err == nil {
		t.Fatalf("missing required environment variable not caught")
	}
}

func TestMissingConfigPath(t *testing.T) {
	os.Unsetenv("CONFIG_PATH")
	if err := requireEnv(); err == nil {
		t.Fatalf("missing required environment variable not caught")
	}
}

func TestMissingWebhookURL(t *testing.T) {
	os.Unsetenv("WEBHOOK_URL")
	if err := requireEnv(); err == nil {
		t.Fatalf("missing required environment variable not caught")
	}
}
