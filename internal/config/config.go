package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// RuntimeConfig represents the structure of the runtime configuration file
type RuntimeConfig struct {
	App string `json:"app"`
}

// ReadRuntimeConfig reads and parses a runtime configuration file
func ReadRuntimeConfig(configPath string) (string, error) {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("error reading config file '%s': %w", configPath, err)
	}

	var config RuntimeConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return "", fmt.Errorf("error parsing config file '%s': %w", configPath, err)
	}

	return config.App, nil
}
