package cmd

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCreate_NoApp(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runCreate(client, createCmdFlags{})
	require.ErrorContains(t, err, "--app flag is required")
}

func TestRunCreate_InvalidEnvVarFormat(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runCreate(client, createCmdFlags{app: "my-app", EnvironmentVariables: []string{"BADFORMAT"}})
	require.ErrorContains(t, err, "invalid environment variable format")
}

func TestRunCreate_InvalidSecretFormat(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runCreate(client, createCmdFlags{app: "my-app", Secrets: []string{"NOSEPARATOR"}})
	require.ErrorContains(t, err, "invalid secret format")
}

func TestRunCreate_APIError(t *testing.T) {
	client := &mockRESTClient{
		putFunc: mockPutError("server error"),
	}
	_, err := runCreate(client, createCmdFlags{app: "my-app", EnvironmentVariables: []string{"K=V"}})
	require.ErrorContains(t, err, "error creating app")
}

func TestRunCreate_Success(t *testing.T) {
	var capturedPath string
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			capturedBody, _ = io.ReadAll(body)
			return json.Unmarshal([]byte(`{"app_url":"https://new-app.example.com"}`), resp)
		},
	}

	appUrl, err := runCreate(client, createCmdFlags{
		app:                  "my-app",
		EnvironmentVariables: []string{"KEY1=val1", "KEY2=val2"},
		Secrets:              []string{"SECRET=sval"},
	})
	require.NoError(t, err)
	assert.Equal(t, "https://new-app.example.com", appUrl)
	assert.Equal(t, "runtime/my-app/deployment", capturedPath)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "val1", req.EnvironmentVariables["KEY1"])
	assert.Equal(t, "sval", req.Secrets["SECRET"])
}

func TestRunCreate_WithRevisionName(t *testing.T) {
	var capturedPath string
	client := &mockRESTClient{
		putFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			return json.Unmarshal([]byte(`{"app_url":"https://app.example.com"}`), resp)
		},
	}

	_, err := runCreate(client, createCmdFlags{app: "my-app", RevisionName: "v2"})
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "revision_name=v2")
}

func TestRunCreate_WithInit(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	client := &mockRESTClient{
		putFunc: mockPutResponse(`{"app_url":"https://init-app.example.com"}`),
	}

	appUrl, err := runCreate(client, createCmdFlags{app: "init-app", Init: true})
	require.NoError(t, err)
	assert.Equal(t, "https://init-app.example.com", appUrl)

	data, err := os.ReadFile("runtime.config.json")
	require.NoError(t, err, "expected runtime.config.json to be created")
	assert.Contains(t, string(data), "init-app")
}

func TestRunCreate_EnvVarWithEqualsInValue(t *testing.T) {
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedBody, _ = io.ReadAll(body)
			return json.Unmarshal([]byte(`{"app_url":"https://app.example.com"}`), resp)
		},
	}

	_, err := runCreate(client, createCmdFlags{app: "my-app", EnvironmentVariables: []string{"KEY=val=with=equals"}})
	require.NoError(t, err)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "val=with=equals", req.EnvironmentVariables["KEY"])
}
