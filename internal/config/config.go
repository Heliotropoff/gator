package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName string = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home_path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base_path := filepath.Join(home_path, configFileName)
	return base_path, nil
}

func Read() (Config, error) {
	configFilepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file_contents, err := os.ReadFile(configFilepath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(file_contents, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil

}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.WriteFile(fullPath, configData, os.FileMode(0644))
	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUsername = userName
	if err := write(*cfg); err != nil {
		return err
	}
	return nil
}
