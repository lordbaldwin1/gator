package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	fp, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fp)
	if err != nil {
		return Config{}, fmt.Errorf("error: failed to open file at path: %s", fp)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUsername = userName
	return write(*cfg)
}

func write(cfg Config) error {
	fp, err := getConfigFilePath()
	if err != nil {
		return errors.New("error: failed to get file path")
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("error: failed to get user's home directory path")
	}

	fp := filepath.Join(homeDir, configFileName)

	return fp, nil
}
