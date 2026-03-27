package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunInit_NoApp(t *testing.T) {
	client := &mockRESTClient{}
	err := runInit(client, initCmdFlags{})
	require.ErrorContains(t, err, "--app flag is required")
}

func TestRunInit_AppNotAccessible(t *testing.T) {
	client := &mockRESTClient{
		getFunc: mockGetError("404 not found"),
	}

	err := runInit(client, initCmdFlags{app: "bad-app"})
	require.ErrorContains(t, err, "does not exist or is not accessible")
}

func TestRunInit_Success_DefaultPath(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	var capturedPath string
	client := &mockRESTClient{
		getFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return json.Unmarshal([]byte(`{"app_url":"https://my-app.example.com"}`), resp)
		},
	}

	err := runInit(client, initCmdFlags{app: "my-app"})
	require.NoError(t, err)
	assert.Equal(t, "runtime/my-app/deployment", capturedPath)

	data, err := os.ReadFile("runtime.config.json")
	require.NoError(t, err, "expected runtime.config.json to be created")
	assert.Contains(t, string(data), "my-app")
}

func TestRunInit_Success_CustomOutPath(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	client := &mockRESTClient{
		getFunc: mockGetResponse(`{"app_url":"https://my-app.example.com"}`),
	}

	outPath := filepath.Join(tmp, "subdir", "custom-config.json")
	err := runInit(client, initCmdFlags{app: "my-app", out: outPath})
	require.NoError(t, err)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err, "expected config file at custom path")
	assert.Contains(t, string(data), "my-app")
}

func TestRunInit_VerifiesAPIBeforeWriting(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	client := &mockRESTClient{
		getFunc: mockGetError("not found"),
	}

	err := runInit(client, initCmdFlags{app: "bad-app"})
	require.Error(t, err)
	require.NoFileExists(t, "runtime.config.json")
}
