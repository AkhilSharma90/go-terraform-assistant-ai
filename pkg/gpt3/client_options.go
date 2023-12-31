package gpt3

import (
	"net/http"
	"time"
)

// ClientOption are options that can be passed when creating a new gpt client.
type ClientOption func(*client) error

//the functions defined below are just helper functions, found from the docs
//they are NOT used, except withHTTPClient which is used in a test file
//having these functions can help serve us later, they don't pollute our current files
//and can exist in separate package and files and can be used if we want to extend functionality

// WithAPIVersion is a client option that allows you to override the default api version of the client.
func WithAPIVersion(apiVersion string) ClientOption {
	return func(c *client) error {
		c.apiVersion = apiVersion
		return nil
	}
}

// WithUserAgent is a client option that allows you to override the default user agent of the client.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithHTTPClient allows you to override the internal http.Client used.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client. The default is 30 seconds. If you are overriding the http client as well, just include
// the timeout there.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}
