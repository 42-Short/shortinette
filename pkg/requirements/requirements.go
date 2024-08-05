package requirements

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/pkg/logger"
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
	var missingValuesString string
	for key, value := range vars {
		if value == "" {
			missingValuesString += "\n" + key
		}
	}
	if missingValuesString != "" {
		return fmt.Errorf("missing environment variables:%s\nSee https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on .env configuration", missingValuesString)
	}
	return nil
}

func ValidateRequirements() error {
	if err := requireEnv(); err != nil {
		return err
	}
	command := "bash"
	args := []string{
		"-c",
		"docker image ls | grep testenv",
	}
	if output, err := testutils.RunCommandLine(".", command, args); err != nil {
		return fmt.Errorf("in order to reduce vulnerability to malicious code, shortinette requires you to have a pre-built Docker image containing all dependencies needed for testing submissions: %s", output)
	}
	logger.Info.Println("all dependencies are already installed")
	return nil
}
