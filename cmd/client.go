package cmd

import "io"

// restClient is the subset of api.RESTClient methods needed by the various commands.
type restClient interface {
	Get(path string, resp interface{}) error
	Delete(path string, resp interface{}) error
	Do(method string, path string, body io.Reader, resp interface{}) error
	Patch(path string, body io.Reader, resp interface{}) error
	Post(path string, body io.Reader, resp interface{}) error
	Put(path string, body io.Reader, resp interface{}) error
}
