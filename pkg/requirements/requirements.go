// requirements provides functions for validating the necessary environment variables
// and dependencies required by the application.
package requirements

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/joho/godotenv"
)

// requireEnv checks for the presence of essential environment variables required by the application.
// It attempts to load these variables from a .env file and validates that they are set.
//
// Returns an error if any required environment variables are missing.
func requireEnv() (err error) {
	if err = godotenv.Load(); err != nil {
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
	if os.Getenv("DEV_MODE") == "" {
		fmt.Printf("DEV_MODE environment variable not set, assuming production")
	}
	if missingValuesString != "" {
		return fmt.Errorf("missing environment variables:%s\nSee https://github.com/42-Short/shortinette/tree/main/.github/docs/DOTENV.md for details on .env configuration", missingValuesString)
	}
	return nil
}

// ValidateRequirements validates the required environment variables and checks for the presence
// of a pre-built Docker image that contains all dependencies needed for testing submissions.
//
// Returns an error if any environment variables are missing or if the Docker image is not found.
func ValidateRequirements() (err error) {
	if err = requireEnv(); err != nil {
		return err
	}
	if err = os.Mkdir("traces/", 0755); err != nil && !os.IsExist(err) {
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
