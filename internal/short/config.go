package short

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Participant struct {
	GithubUserName string
	IntraLogin     string
}

type Config struct {
	StartDate    time.Time
	EndDate      time.Time
	Participants []Participant
}

func getConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable not set")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	}

	var rawConfig struct {
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		Participants []struct {
			GithubUserName string `json:"github_username"`
			IntraLogin     string `json:"intra_login"`
		} `json:"participants"`
	}

	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}

	startDate, err := time.Parse("02.01.2006", rawConfig.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("02.01.2006", rawConfig.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	var participants []Participant
	for _, p := range rawConfig.Participants {
		participants = append(participants, Participant{
			GithubUserName: p.GithubUserName,
			IntraLogin:     p.IntraLogin,
		})
	}

	config := &Config{
		StartDate:    startDate,
		EndDate:      endDate,
		Participants: participants,
	}

	return config, nil
}
