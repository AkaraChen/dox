package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads a do.yaml configuration file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
	 return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
	 return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate version
	if config.Version != 1 && config.Version != 0 {
	 return nil, fmt.Errorf("unsupported config version: %d", config.Version)
	}

	return &config, nil
}

// FindConfigFile looks for do.yaml in the specified directory
func FindConfigFile(dir string) (string, error) {
	configPath := filepath.Join(dir, "do.yaml")
	if _, err := os.Stat(configPath); err == nil {
	 return configPath, nil
	}
	return "", fmt.Errorf("no do.yaml found in %s", dir)
}

// LoadConfigFromDirectory loads do.yaml from directory if it exists
func LoadConfigFromDirectory(dir string) (*Config, string, error) {
	configPath, err := FindConfigFile(dir)
	if err != nil {
	 // No config file is not an error
	 return nil, "", nil
	}

	config, err := LoadConfig(configPath)
	if err != nil {
	 return nil, "", err
	}

	return config, configPath, nil
}
