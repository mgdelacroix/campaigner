package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

type Config struct {
	GithubToken string `json:"github_token"`
	JiraToken   string `json:"jira_token"`
}

func getConfigPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.campaigner", nil
}

func ReadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); err != nil {
		return &Config{}, nil 
	}

	fileContents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("there was a problem reading the config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(fileContents, &config); err != nil {
		return nil, fmt.Errorf("there was a problem parsing the config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	marshaledConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configPath, marshaledConfig, 0600); err != nil {
		return fmt.Errorf("cannot save the config: %w", err)
	}
	return nil
}
