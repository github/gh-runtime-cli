package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunDeploy_NoDir(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	client := &mockRESTClient{}
	err = runDeploy(client, deployCmdFlags{app: "my-app"})
	require.ErrorContains(t, err, "--dir flag is required")
}

func TestRunDeploy_NoAppName(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	deployDir := filepath.Join(tmp, "dist")
	require.NoError(t, os.MkdirAll(deployDir, 0755))

	client := &mockRESTClient{}
	err = runDeploy(client, deployCmdFlags{dir: deployDir})
	require.ErrorContains(t, err, "--app flag is required")
}

func TestRunDeploy_DirNotExist(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	client := &mockRESTClient{}
	err = runDeploy(client, deployCmdFlags{dir: "/nonexistent/path", app: "my-app"})
	require.ErrorContains(t, err, "does not exist")
}

func TestRunDeploy_Success(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	deployDir := filepath.Join(tmp, "dist")
	require.NoError(t, os.MkdirAll(deployDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(deployDir, "index.html"), []byte("<html></html>"), 0644))

	var capturedPath string
	client := &mockRESTClient{
		postFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			return nil
		},
	}

	err = runDeploy(client, deployCmdFlags{dir: deployDir, app: "my-app"})
	require.NoError(t, err)
	assert.Equal(t, "runtime/my-app/deployment/bundle", capturedPath)
	require.NoFileExists(t, deployDir+".zip")
}

func TestRunDeploy_WithRevisionAndSha(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	deployDir := filepath.Join(tmp, "dist")
	require.NoError(t, os.MkdirAll(deployDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(deployDir, "app.js"), []byte("console.log('hi')"), 0644))

	var capturedPath string
	client := &mockRESTClient{
		postFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			return nil
		},
	}

	err = runDeploy(client, deployCmdFlags{dir: deployDir, app: "my-app", revisionName: "v2", sha: "abc123"})
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "revision_name=v2")
	assert.Contains(t, capturedPath, "revision=abc123")
}

func TestRunDeploy_APIError(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	deployDir := filepath.Join(tmp, "dist")
	require.NoError(t, os.MkdirAll(deployDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(deployDir, "index.html"), []byte("<html></html>"), 0644))

	client := &mockRESTClient{
		postFunc: mockPostError("upload failed"),
	}

	err = runDeploy(client, deployCmdFlags{dir: deployDir, app: "my-app"})
	require.ErrorContains(t, err, "error deploying app")
}

func TestRunDeploy_WithConfigFile(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	deployDir := filepath.Join(tmp, "dist")
	require.NoError(t, os.MkdirAll(deployDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(deployDir, "index.html"), []byte("<html></html>"), 0644))

	configPath := filepath.Join(tmp, "my-config.json")
	require.NoError(t, os.WriteFile(configPath, []byte(`{"app":"config-deploy-app"}`), 0644))

	var capturedPath string
	client := &mockRESTClient{
		postFunc: func(path string, body io.Reader, resp interface{}) error {
			capturedPath = path
			return nil
		},
	}

	err = runDeploy(client, deployCmdFlags{dir: deployDir, config: configPath})
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "config-deploy-app")
}
