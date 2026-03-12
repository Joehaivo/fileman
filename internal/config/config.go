package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	UseEnglish bool `json:"use_english"`
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "fileman", "config.json")
	}

	return filepath.Join(homeDir, ".config", "fileman", "config.json")
}

func LoadConfig() *Config {
	cfg := &Config{UseEnglish: false}

	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(data, cfg)
	return cfg
}

func SaveConfig(cfg *Config) error {
	path := GetConfigPath()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
