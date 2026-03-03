package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	NotesDir string `json:"notes_dir"`
}

func LoadConfig() (*Config, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	cfgPath := filepath.Join(cfgDir, "silo", "config.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ConfigExists() bool {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return false
	}

	cfgPath := filepath.Join(cfgDir, "silo", "config.json")
	_, err = os.Stat(cfgPath)
	return err == nil
}

func SaveConfig(cfg *Config) error {
	//* Unix/Linux: $XDG_CONFIG_HOME or $HOME/.config
	//* Darwin    : $HOME/Library/Application Support
	//* Windows   : %AppData%
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	//? userConfigDir/silo
	siloCfgDir := filepath.Join(userCfgDir, "silo")
	if err := os.MkdirAll(siloCfgDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	siloCfgPath := filepath.Join(siloCfgDir, "config.json")
	return os.WriteFile(siloCfgPath, data, 0644)
}

func GetConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(dir, "silo", "config.json"), nil
}

