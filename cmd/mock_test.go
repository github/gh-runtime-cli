package cmd

import (
	"encoding/json"
	"fmt"
	"io"
)

// mockRESTClient implements the restClient interface for testing.
type mockRESTClient struct {
	getFunc    func(path string, resp interface{}) error
	putFunc    func(path string, body io.Reader, resp interface{}) error
	deleteFunc func(path string, resp interface{}) error
	postFunc   func(path string, body io.Reader, resp interface{}) error
	patchFunc  func(path string, body io.Reader, resp interface{}) error
	doFunc     func(method string, path string, body io.Reader, resp interface{}) error
}

func (m *mockRESTClient) Get(path string, resp interface{}) error {
	if m.getFunc != nil {
		return m.getFunc(path, resp)
	}
	return fmt.Errorf("Get not implemented")
}

func (m *mockRESTClient) Put(path string, body io.Reader, resp interface{}) error {
	if m.putFunc != nil {
		return m.putFunc(path, body, resp)
	}
	return fmt.Errorf("Put not implemented")
}

func (m *mockRESTClient) Delete(path string, resp interface{}) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(path, resp)
	}
	return fmt.Errorf("Delete not implemented")
}

func (m *mockRESTClient) Post(path string, body io.Reader, resp interface{}) error {
	if m.postFunc != nil {
		return m.postFunc(path, body, resp)
	}
	return fmt.Errorf("Post not implemented")
}

func (m *mockRESTClient) Patch(path string, body io.Reader, resp interface{}) error {
	if m.patchFunc != nil {
		return m.patchFunc(path, body, resp)
	}
	return fmt.Errorf("Patch not implemented")
}

func (m *mockRESTClient) Do(method string, path string, body io.Reader, resp interface{}) error {
	if m.doFunc != nil {
		return m.doFunc(method, path, body, resp)
	}
	return fmt.Errorf("Do not implemented")
}

// mockGetResponse is a helper that configures the mock to return a JSON-decoded response.
func mockGetResponse(jsonBody string) func(path string, resp interface{}) error {
	return func(path string, resp interface{}) error {
		return json.Unmarshal([]byte(jsonBody), resp)
	}
}

// mockGetError is a helper that configures the mock to return an error.
func mockGetError(errMsg string) func(path string, resp interface{}) error {
	return func(path string, resp interface{}) error {
		return fmt.Errorf("%s", errMsg)
	}
}

// mockPutResponse is a helper that configures the mock to return a JSON-decoded response for Put.
func mockPutResponse(jsonBody string) func(path string, body io.Reader, resp interface{}) error {
	return func(path string, body io.Reader, resp interface{}) error {
		return json.Unmarshal([]byte(jsonBody), resp)
	}
}

// mockPutError is a helper that configures the mock to return an error for Put.
func mockPutError(errMsg string) func(path string, body io.Reader, resp interface{}) error {
	return func(path string, body io.Reader, resp interface{}) error {
		return fmt.Errorf("%s", errMsg)
	}
}

// mockDeleteResponse is a helper for Delete.
func mockDeleteResponse(jsonBody string) func(path string, resp interface{}) error {
	return func(path string, resp interface{}) error {
		return json.Unmarshal([]byte(jsonBody), resp)
	}
}

// mockDeleteError is a helper that configures the mock to return an error for Delete.
func mockDeleteError(errMsg string) func(path string, resp interface{}) error {
	return func(path string, resp interface{}) error {
		return fmt.Errorf("%s", errMsg)
	}
}

// mockPostSuccess is a helper for Post returning nil.
func mockPostSuccess() func(path string, body io.Reader, resp interface{}) error {
	return func(path string, body io.Reader, resp interface{}) error {
		return nil
	}
}

// mockPostError is a helper that configures the mock to return an error for Post.
func mockPostError(errMsg string) func(path string, body io.Reader, resp interface{}) error {
	return func(path string, body io.Reader, resp interface{}) error {
		return fmt.Errorf("%s", errMsg)
	}
}
