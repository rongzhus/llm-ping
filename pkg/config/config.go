package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	ConfigFileName = "llm-ping.json"
)

type ModelConfig struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	Endpoint  string `json:"endpoint"`
	APIKey    string `json:"apiKey,omitempty"`
	ModelName string `json:"modelName"`
}

type Config struct {
	Models []ModelConfig `json:"models"`
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".llm-ping", ConfigFileName), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Models: []ModelConfig{}}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func AddModel(model ModelConfig) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	for i, m := range cfg.Models {
		if m.Name == model.Name {
			cfg.Models[i] = model
			return SaveConfig(cfg)
		}
	}

	cfg.Models = append(cfg.Models, model)
	return SaveConfig(cfg)
}

func RemoveModel(name string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	var newModels []ModelConfig
	for _, m := range cfg.Models {
		if m.Name != name {
			newModels = append(newModels, m)
		}
	}

	cfg.Models = newModels
	return SaveConfig(cfg)
}

func GetModel(name string) (*ModelConfig, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	for _, m := range cfg.Models {
		if m.Name == name {
			return &m, nil
		}
	}

	return nil, nil
}

func ListModels() ([]ModelConfig, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Models, nil
}
