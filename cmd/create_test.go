package cmd

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildCreateResponse(r createResp, resp interface{}) {
	if r.AppUrl == "" {
		r.AppUrl = "https://test-app.example.com"
	}
	if r.ID == "" {
		r.ID = "test-app-id"
	}
	*resp.(*createResp) = r
}

func TestRunCreate_NoAppOrName(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runCreate(client, createCmdFlags{})
	require.ErrorContains(t, err, "either --app or --name flag is required")
}

func TestRunCreate_InvalidEnvVarFormat(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runCreate(client, createCmdFlags{app: "my-app", environmentVariables: []string{"BADFORMAT"}})
	require.ErrorContains(t, err, "invalid environment variable format")
}

func TestRunCreate_InvalidSecretFormat(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runCreate(client, createCmdFlags{app: "my-app", secrets: []string{"NOSEPARATOR"}})
	require.ErrorContains(t, err, "invalid secret format")
}

func TestRunCreate_APIError(t *testing.T) {
	client := &mockRESTClient{
		putFunc: mockPutError("server error"),
	}
	_, err := runCreate(client, createCmdFlags{app: "my-app", environmentVariables: []string{"K=V"}})
	require.ErrorContains(t, err, "error creating app")
}

func TestRunCreate_Success(t *testing.T) {
	var capturedPath string
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			capturedBody, _ = io.ReadAll(body)
			buildCreateResponse(createResp{AppUrl: "https://my-app.example.com"}, resp)
			return nil
		},
	}

	resp, err := runCreate(client, createCmdFlags{
		app:                  "my-app",
		environmentVariables: []string{"KEY1=val1", "KEY2=val2"},
		secrets:              []string{"SECRET=sval"},
	})
	require.NoError(t, err)
	assert.Equal(t, "https://my-app.example.com", resp.AppUrl)
	assert.Equal(t, "runtime/my-app/deployment", capturedPath)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "val1", req.EnvironmentVariables["KEY1"])
	assert.Equal(t, "sval", req.Secrets["SECRET"])
}

func TestRunCreate_WithRevisionName(t *testing.T) {
	var capturedPath string
	client := &mockRESTClient{
		putFunc: func(path string, _ io.Reader, resp interface{}) error {
			capturedPath = path
			buildCreateResponse(createResp{AppUrl: "https://my-app.example.com"}, resp)
			return nil
		},
	}

	_, err := runCreate(client, createCmdFlags{app: "my-app", revisionName: "v2"})
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "revision_name=v2")
}

func TestRunCreate_WithInit(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	client := &mockRESTClient{
		putFunc: func(_ string, _ io.Reader, resp interface{}) error {
			buildCreateResponse(createResp{AppUrl: "https://init-app.example.com", ID: "init-app-id"}, resp)
			return nil
		},
	}

	resp, err := runCreate(client, createCmdFlags{app: "init-app", init: true})
	require.NoError(t, err)
	assert.Equal(t, "https://init-app.example.com", resp.AppUrl)

	data, err := os.ReadFile("runtime.config.json")
	require.NoError(t, err, "expected runtime.config.json to be created")
	assert.Contains(t, string(data), "init-app-id")
}

func TestRunCreate_EnvVarWithEqualsInValue(t *testing.T) {
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(_ string, body io.Reader, resp interface{}) error {
			capturedBody, _ = io.ReadAll(body)
			buildCreateResponse(createResp{AppUrl: "https://my-app.example.com"}, resp)
			return nil
		},
	}

	_, err := runCreate(client, createCmdFlags{app: "my-app", environmentVariables: []string{"KEY=val=with=equals"}})
	require.NoError(t, err)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "val=with=equals", req.EnvironmentVariables["KEY"])
}

func TestRunCreate_WithName(t *testing.T) {
	var capturedPath string
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			capturedBody, _ = io.ReadAll(body)
			buildCreateResponse(createResp{AppUrl: "https://my-new-app.example.com", ID: "abc-123"}, resp)
			return nil
		},
	}

	resp, err := runCreate(client, createCmdFlags{name: "my-new-app"})
	require.NoError(t, err)
	assert.Equal(t, "https://my-new-app.example.com", resp.AppUrl)
	assert.Equal(t, "abc-123", resp.ID)
	assert.Equal(t, "runtime", capturedPath)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "my-new-app", req.Name)
}

func TestRunCreate_WithNameAndApp(t *testing.T) {
	var capturedPath string
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			capturedBody, _ = io.ReadAll(body)
			buildCreateResponse(createResp{AppUrl: "https://my-app.example.com", ID: "app-123"}, resp)
			return nil
		},
	}

	resp, err := runCreate(client, createCmdFlags{app: "my-app", name: "my-new-name"})
	require.NoError(t, err)
	assert.Equal(t, "https://my-app.example.com", resp.AppUrl)
	assert.Equal(t, "runtime/my-app/deployment", capturedPath)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "my-new-name", req.Name)
}

func TestRunCreate_WithNameAndInit(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(_ string, body io.Reader, resp interface{}) error {
			capturedBody, _ = io.ReadAll(body)
			buildCreateResponse(createResp{AppUrl: "https://named-app.example.com", ID: "def-456"}, resp)
			return nil
		},
	}

	_, err := runCreate(client, createCmdFlags{name: "named-app", init: true})
	require.NoError(t, err)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "named-app", req.Name)

	data, err := os.ReadFile("runtime.config.json")
	require.NoError(t, err)
	assert.Contains(t, string(data), "def-456")
}

func TestRunCreate_WithVisibility(t *testing.T) {
	var capturedBody []byte
	client := &mockRESTClient{
		putFunc: func(_ string, body io.Reader, resp interface{}) error {
			capturedBody, _ = io.ReadAll(body)
			buildCreateResponse(createResp{AppUrl: "https://my-app.example.com"}, resp)
			return nil
		},
	}

	_, err := runCreate(client, createCmdFlags{app: "my-app", visibility: "github"})
	require.NoError(t, err)

	var req createReq
	json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "github", req.Visibility)
}

func TestRunCreate_ResponseWithID(t *testing.T) {
	client := &mockRESTClient{
		putFunc: func(_ string, _ io.Reader, resp interface{}) error {
			buildCreateResponse(createResp{AppUrl: "https://my-app.example.com", ID: "xyz-789"}, resp)
			return nil
		},
	}

	resp, err := runCreate(client, createCmdFlags{app: "my-app"})
	require.NoError(t, err)
	assert.Equal(t, "https://my-app.example.com", resp.AppUrl)
	assert.Equal(t, "xyz-789", resp.ID)
}

func TestRunCreate_InitWithoutIDInResponse(t *testing.T) {
	client := &mockRESTClient{
		putFunc: func(_ string, _ io.Reader, resp interface{}) error {
			b, _ := json.Marshal(createResp{AppUrl: "https://my-app.example.com"})
			return json.Unmarshal(b, resp)
		},
	}

	_, err := runCreate(client, createCmdFlags{app: "my-app", init: true})
	require.ErrorContains(t, err, "server did not return an app ID")
}
