package gpt3_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 net/http.RoundTripper

// TestInitNewClient tests the scenario where a new GPT-3 client is initialized with a test key.
func TestInitNewClient(t *testing.T) {
	client := gpt3.NewClient("test-key")
	assert.NotNil(t, client)
}

func fakeHTTPClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}
	return rt, &http.Client{
		Transport: rt,
	}
}

// TestRequestCreationFails tests the scenario where request creation fails for various API calls.
// It sets up a fake HTTP client and mocks the round trip to return a request error.
// Then, it defines a list of test cases, each representing an API call and the expected error string.
// For each test case, it calls the API function and asserts that the returned error matches the expected error string.
// Finally, it asserts that the response is nil for each test case.
func TestRequestCreationFails(t *testing.T) {
	ctx := context.Background()
	rt, httpClient := fakeHTTPClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	rt.RoundTripReturns(nil, errors.New("request error"))

	type testCase struct {
		name        string
		apiCall     func() (interface{}, error)
		errorString string
	}

	testCases := []testCase{
		{
			"Engines",
			func() (interface{}, error) {
				return client.Engines(ctx)
			},
			"Get \"https://api.openai.com/v1/engines\": request error",
		},
		{
			"Engine",
			func() (interface{}, error) {
				return client.Engine(ctx, gpt3.DefaultEngine)
			},
			"Get \"https://api.openai.com/v1/engines/davinci\": request error",
		},
		{
			"Completion",
			func() (interface{}, error) {
				return client.Completion(ctx, gpt3.CompletionRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/davinci/completions\": request error",
		},
		{
			"CompletionStream",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStream(ctx, gpt3.CompletionRequest{}, onData)
			},
			"Post \"https://api.openai.com/v1/engines/davinci/completions\": request error",
		},
		{
			"CompletionWithEngine",
			func() (interface{}, error) {
				return client.CompletionWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/ada/completions\": request error",
		},
		{
			"CompletionStreamWithEngine",
			func() (interface{}, error) {
				var rsp *gpt3.CompletionResponse
				onData := func(data *gpt3.CompletionResponse) {
					rsp = data
				}
				return rsp, client.CompletionStreamWithEngine(ctx, gpt3.AdaEngine, gpt3.CompletionRequest{}, onData)
			},
			"Post \"https://api.openai.com/v1/engines/ada/completions\": request error",
		},
		{
			"Edits",
			func() (interface{}, error) {
				return client.Edits(ctx, gpt3.EditsRequest{})
			},
			"Post \"https://api.openai.com/v1/edits\": request error",
		},
		{
			"Search",
			func() (interface{}, error) {
				return client.Search(ctx, gpt3.SearchRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/davinci/search\": request error",
		},
		{
			"SearchWithEngine",
			func() (interface{}, error) {
				return client.SearchWithEngine(ctx, gpt3.AdaEngine, gpt3.SearchRequest{})
			},
			"Post \"https://api.openai.com/v1/engines/ada/search\": request error",
		},
		{
			"Embeddings",
			func() (interface{}, error) {
				return client.Embeddings(ctx, gpt3.EmbeddingsRequest{})
			},
			"Post \"https://api.openai.com/v1/embeddings\": request error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rsp, err := tc.apiCall()
			assert.EqualError(t, err, tc.errorString)
			assert.Nil(t, rsp)
		})
	}
}

type errReader int

// Read is a method that implements the Read function of the io.Reader interface.
func (errReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestResponses(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Create a fake HTTP client for testing
	rt, httpClient := fakeHTTPClient()

	// Create a new GPT-3 client with a test key and the fake HTTP client
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))

	// Define a test case struct to hold the name, API call function, and expected response object
	type testCase struct {
		name           string
		apiCall        func() (interface{}, error)
		responseObject interface{}
	}

	// Define a list of test cases
	testCases := []testCase{
		{
			name: "Engines",
			apiCall: func() (interface{}, error) {
				return client.Engines(ctx)
			},
			responseObject: &gpt3.EnginesResponse{
				Data: []gpt3.EngineObject{
					{
						ID:     "123",
						Object: "list",
						Owner:  "owner",
						Ready:  true,
					},
				},
			},
		},
		// Add more test cases...
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("bad status codes", func(t *testing.T) {
				// Test different HTTP status codes
				for _, code := range []int{400, 401, 404, 422, 500} {
					// Mock the HTTP response with the specified status code and an error body
					mockResponse := &http.Response{
						StatusCode: code,
						Body:       io.NopCloser(errReader(0)),
					}

					// Set the mock response for the fake HTTP client
					rt.RoundTripReturns(mockResponse, nil)

					// Make the API call and assert the expected error
					rsp, err := tc.apiCall()
					assert.Nil(t, rsp)
					assert.EqualError(t, err, "failed to read from body: read error")

					// Mock the HTTP response with the specified status code and an unknown error string
					mockResponse = &http.Response{
						StatusCode: code,
						Body:       io.NopCloser(bytes.NewBufferString("unknown error")),
					}

					// Set the mock response for the fake HTTP client
					rt.RoundTripReturns(mockResponse, nil)

					// Make the API call and assert the expected error
					rsp, err = tc.apiCall()
					assert.Nil(t, rsp)
					assert.EqualError(t, err, fmt.Sprintf("[%d:Unexpected] unknown error", code))

					// Mock the HTTP response with the specified status code and a JSON APIErrorResponse
					apiErrorResponse := &gpt3.APIErrorResponse{
						Error: gpt3.APIError{
							Type:    "test_type",
							Message: "test message",
						},
					}

					data, err := json.Marshal(apiErrorResponse)
					assert.NoError(t, err)

					mockResponse = &http.Response{
						StatusCode: code,
						Body:       io.NopCloser(bytes.NewBuffer(data)),
					}

					// Set the mock response for the fake HTTP client
					rt.RoundTripReturns(mockResponse, nil)

					// Make the API call and assert the expected error
					rsp, err = tc.apiCall()
					assert.Nil(t, rsp)
					assert.EqualError(t, err, fmt.Sprintf("[%d:test_type] test message", code))
					apiErrorResponse.Error.StatusCode = code
					assert.Equal(t, apiErrorResponse.Error, err)
				}
			})

			t.Run("success code json decode failure", func(t *testing.T) {
				// Mock the HTTP response with a success status code and an invalid JSON body
				mockResponse := &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
				}

				// Set the mock response for the fake HTTP client
				rt.RoundTripReturns(mockResponse, nil)

				// Make the API call and assert the expected error
				rsp, err := tc.apiCall()
				assert.Error(t, err, "invalid json response: invalid character 'i' looking for beginning of value")
				assert.Nil(t, rsp)
			})

			// Skip streaming/nil response objects here as those will be tested separately
			if tc.responseObject != nil {
				t.Run("successful response", func(t *testing.T) {
					// Marshal the expected response object to JSON
					data, err := json.Marshal(tc.responseObject)
					assert.NoError(t, err)

					// Mock the HTTP response with a success status code and the JSON body
					mockResponse := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBuffer(data)),
					}

					// Set the mock response for the fake HTTP client
					rt.RoundTripReturns(mockResponse, nil)

					// Make the API call and assert the expected response
					rsp, err := tc.apiCall()
					assert.NoError(t, err)
					assert.Equal(t, tc.responseObject, rsp)
				})
			}
		})
	}
}

// TODO: add streaming response tests
