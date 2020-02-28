package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

func getConfigPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.campaigner", nil
}

func ReadConfig() (*model.Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); err != nil {
		return &model.Config{}, nil
	}

	fileContents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("there was a problem reading the config file: %w", err)
	}

	var config model.Config
	if err := json.Unmarshal(fileContents, &config); err != nil {
		return nil, fmt.Errorf("there was a problem parsing the config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *model.Config) error {
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
