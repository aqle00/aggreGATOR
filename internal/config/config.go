package config

import (
	"encoding/json"

	"os"
	"path/filepath"
)

// path to config file
const configFileName = ".gatorconfig.json"

// struct to hold config data
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// get path to config file
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homePath, configFileName)
	return path, nil
}

// reads the config file, return a config struct + error
func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var configs Config
	dec := json.NewDecoder(file)
	if err := dec.Decode(&configs); err != nil {
		return Config{}, err
	}

	return configs, nil
}

// Set username, write to config file
func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}

// Set DBURL, write to config file
func (cfg *Config) SetDBURL(dbURL string) error {
	cfg.DBURL = dbURL
	return write(*cfg)
}

// write to config file
func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(cfg); err != nil {
		return err
	}

	return nil
}
