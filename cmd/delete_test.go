package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunDelete_NoApp(t *testing.T) {
	client := &mockRESTClient{}
	_, err := runDelete(client, deleteCmdFlags{})
	require.ErrorContains(t, err, "--app flag is required")
}

func TestRunDelete_Success(t *testing.T) {
	var capturedPath string
	client := &mockRESTClient{
		deleteFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return nil
		},
	}

	response, err := runDelete(client, deleteCmdFlags{app: "my-app"})
	require.NoError(t, err)
	assert.Equal(t, "runtime/my-app/deployment", capturedPath)
	assert.Equal(t, "my-app", response)
}

func TestRunDelete_WithRevisionName(t *testing.T) {
	var capturedPath string
	client := &mockRESTClient{
		deleteFunc: func(path string, resp interface{}) error {
			capturedPath = path
			return nil
		},
	}

	_, err := runDelete(client, deleteCmdFlags{app: "my-app", revisionName: "v2"})
	require.NoError(t, err)
	assert.Contains(t, capturedPath, "revision_name=v2")
}

func TestRunDelete_APIError(t *testing.T) {
	client := &mockRESTClient{
		deleteFunc: mockDeleteError("not found"),
	}

	_, err := runDelete(client, deleteCmdFlags{app: "my-app"})
	require.ErrorContains(t, err, "error deleting app")
}
