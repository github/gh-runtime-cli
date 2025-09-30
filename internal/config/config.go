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

// ResolveAppName resolves the app name using the priority order:
// 1. appFlag (--app) if provided
// 2. configPath (--config) if provided
// 3. runtime.config.json in current directory if it exists
// Returns an error if no app name can be resolved
func ResolveAppName(appFlag, configPath string) (string, error) {
	// Priority 1: Use --app flag if provided
	if appFlag != "" {
		return appFlag, nil
	}

	// Priority 2: Use --config file if provided
	if configPath != "" {
		return ReadRuntimeConfig(configPath)
	}

	// Priority 3: Try default runtime.config.json
	if _, err := os.Stat("runtime.config.json"); err == nil {
		appName, err := ReadRuntimeConfig("runtime.config.json")
		if err != nil {
			return "", fmt.Errorf("found runtime.config.json but failed to read it: %v", err)
		}
		return appName, nil
	}

	// No app name could be resolved
	return "", fmt.Errorf("--app flag is required, --config must be specified, or runtime.config.json must exist in current directory")
}
