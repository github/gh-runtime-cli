package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunGet_NoAppName(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	client := &mockRESTClient{}
	_, err = runGet(client, getCmdFlags{})
	require.ErrorContains(t, err, "--app flag is required")
}

func TestRunGet_Success(t *testing.T) {
	var capturedPath string
	client := &mockRESTClient{
		getFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return json.Unmarshal([]byte(`{"app_url":"https://my-app.example.com"}`), resp)
		},
	}

	appUrl, err := runGet(client, getCmdFlags{app: "my-app"})
	require.NoError(t, err)
	assert.Equal(t, "runtime/my-app/deployment", capturedPath)
	assert.Equal(t, "https://my-app.example.com", appUrl)
}

func TestRunGet_WithRevisionName(t *testing.T) {
	var capturedPath string
	client := &mockRESTClient{
		getFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return json.Unmarshal([]byte(`{"app_url":"https://my-app-v2.example.com"}`), resp)
		},
	}

	appUrl, err := runGet(client, getCmdFlags{app: "my-app", revisionName: "v2"})
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "runtime/my-app/deployment")
	assert.Contains(t, capturedPath, "revision_name=v2")
	assert.Equal(t, "https://my-app-v2.example.com", appUrl)
}

func TestRunGet_APIError(t *testing.T) {
	client := &mockRESTClient{
		getFunc: mockGetError("server error 500"),
	}

	_, err := runGet(client, getCmdFlags{app: "my-app"})
	require.ErrorContains(t, err, "retrieving app details")
}

func TestRunGet_WithConfigFile(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	configPath := filepath.Join(tmp, "my-config.json")
	require.NoError(t, os.WriteFile(configPath, []byte(`{"app":"config-app"}`), 0644))

	var capturedPath string
	client := &mockRESTClient{
		getFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return json.Unmarshal([]byte(`{"app_url":"https://config-app.example.com"}`), resp)
		},
	}

	appUrl, err := runGet(client, getCmdFlags{config: configPath})
	require.NoError(t, err)
	assert.Equal(t, capturedPath, "runtime/config-app/deployment")
	assert.Equal(t, "https://config-app.example.com", appUrl)
}

func TestRunGet_DefaultConfigFile(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	require.NoError(t, os.WriteFile(filepath.Join(tmp, "runtime.config.json"), []byte(`{"app":"default-app"}`), 0644))

	var capturedPath string
	client := &mockRESTClient{
		getFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return json.Unmarshal([]byte(`{"app_url":"https://default-app.example.com"}`), resp)
		},
	}

	appUrl, err := runGet(client, getCmdFlags{})
	require.NoError(t, err)
	assert.Equal(t, capturedPath, "runtime/default-app/deployment")
	assert.Equal(t, "https://default-app.example.com", appUrl)
}
