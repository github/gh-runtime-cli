package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// runtimeConfig represents the structure of the runtime configuration file
type runtimeConfig struct {
	App string `json:"app"`
}

// readRuntimeConfig reads and parses a runtime configuration file,
// returning the app name specified in the config
func readRuntimeConfig(configPath string) (string, error) {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("error reading config file '%s': %w", configPath, err)
	}

	var config runtimeConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return "", fmt.Errorf("error parsing config file '%s': %w", configPath, err)
	}

	return config.App, nil
}
