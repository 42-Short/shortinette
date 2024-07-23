package requirements

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/joho/godotenv"
)

func requireEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("could not load .env file: %v", err)
	}
	vars := map[string]string{
		"GITHUB_ADMIN":        os.Getenv("GITHUB_ADMIN"),
		"GITHUB_EMAIL":        os.Getenv("GITHUB_EMAIL"),
		"GITHUB_TOKEN":        os.Getenv("GITHUB_TOKEN"),
		"GITHUB_ORGANISATION": os.Getenv("GITHUB_ORGANISATION"),
		"CONFIG_PATH":         os.Getenv("CONFIG_PATH"),
		"WEBHOOK_URL":         os.Getenv("WEBHOOK_URL"),
	}
	for key, value := range vars {
		if value == "" {
			return fmt.Errorf("%s environment variable not set", key)
		}
	}
	return nil
}

func ValidateRequirements() error {
	if err := requireEnv(); err != nil {
		return err
	}
	if output, err := testutils.RunExecutable("./scripts/startup.sh"); err != nil {
		return fmt.Errorf(output)
	}
	logger.Info.Println("all dependencies are already installed")
	return nil
}
