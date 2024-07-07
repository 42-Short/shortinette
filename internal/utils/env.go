package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func RequireEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("could not load .env file: %v", err)
	}
	vars := map[string]string{
		"GITHUB_USER":         os.Getenv("GITHUB_USER"),
		"GITHUB_TOKEN":        os.Getenv("GITHUB_TOKEN"),
		"GITHUB_ORGANISATION": os.Getenv("GITHUB_ORGANISATION"),
		"CONFIG_PATH":         os.Getenv("CONFIG_PATH"),
	}
	for key, value := range vars {
		if value == "" {
			return fmt.Errorf("%s environment variable not set", key)
		}
	}
	return nil
}
