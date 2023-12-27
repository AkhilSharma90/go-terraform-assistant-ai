package gpt3

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultAPIVersion     = "2023-03-15-preview"
	defaultUserAgent      = "kubectl-openai"
	defaultTimeoutSeconds = 30
)

// A Client is an API client to communicate with the OpenAI gpt-3 APIs.
type Client interface {
	// ChatCompletion creates a completion with the Chat completion endpoint which
	// is what powers the ChatGPT experience.
	ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error)

	// Completion creates a completion with the default engine. This is the main endpoint of the API
	// which auto-completes based on the given prompt.
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)

	// CompletionStream creates a completion with the default engine and streams the results through
	// multiple calls to onData.
	CompletionStream(ctx context.Context, request CompletionRequest, onData func(*CompletionResponse)) error

	// Given a prompt and an instruction, the model will return an edited version of the prompt.
	Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error)

	// Search performs a semantic search over a list of documents with the default engine.
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)

	// Returns an embedding using the provided request.
	Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error)
}

type client struct {
	endpoint       string
	apiKey         string
	deploymentName string
	apiVersion     string
	userAgent      string
	httpClient     *http.Client
}

//Getting called in the NewOAIClients function in completion.go file
// NewClient returns a new OpenAI GPT-3 API client. An apiKey is required to use the client.
// NewClient creates a new GPT-3 client with the specified endpoint, API key, deployment name, and optional client options.
func NewClient(endpoint string, apiKey string, deploymentName string, options ...ClientOption) (Client, error) {
	// Create a new HTTP client with a default timeout.
	httpClient := &http.Client{
		Timeout: defaultTimeoutSeconds * time.Second,
	}

	// Create a new client instance with the provided parameters.
	c := &client{
		endpoint:       endpoint,
		apiKey:         apiKey,
		deploymentName: deploymentName,
		apiVersion:     defaultAPIVersion,
		userAgent:      defaultUserAgent,
		httpClient:     httpClient,
	}

	// Apply any additional client options provided.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

//Getting called in openai.go file, azureGPTCompletion function
// Completion sends a completion request to the OpenAI API and returns the completion response.
func (c *client) Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	// Set the Stream field of the request to false
	request.Stream = false

	// Create a new request using the context, HTTP method, and endpoint URL
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/openai/deployments/%s/completions", c.deploymentName), request)
	if err != nil {
		return nil, err
	}

	// Perform the request and get the response
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	// Create a new CompletionResponse object to store the response data
	output := new(CompletionResponse)

	// Parse the response and populate the output object
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}

	// Return the output object and nil error
	return output, nil
}

//Getting called in openai.go file, azureGPTChatCompletion function
// ChatCompletion sends a chat completion request to the OpenAI API and returns the response.
func (c *client) ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	request.Stream = false

	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/openai/deployments/%s/chat/completions", c.deploymentName), request)
	if err != nil {
		return nil, err
	}

	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ChatCompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

var (
	dataPrefix   = []byte("data: ")
	doneSequence = []byte("[DONE]")
)

//FUNCTION NOT GETTING USED
// CompletionStream is a method that allows streaming of completion responses from the OpenAI API.
func (c *client) CompletionStream(ctx context.Context, request CompletionRequest, onData func(*CompletionResponse)) error {
	// Set the stream flag to true in the request
	request.Stream = true

	// Create a new request using the provided context, HTTP method, and endpoint
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/openai/deployments/%s/completions", c.deploymentName), request)
	if err != nil {
		return err
	}

	// Perform the request and get the response
	resp, err := c.performRequest(req)
	if err != nil {
		return err
	}

	// Create a new reader to read the response body
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	// Read the response body line by line
	for {
		// Read a line from the response body
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		// Trim any extra whitespace from the line
		line = bytes.TrimSpace(line)

		// Check if the line is a data event
		if !bytes.HasPrefix(line, dataPrefix) {
			continue
		}

		// Remove the data prefix from the line
		line = bytes.TrimPrefix(line, dataPrefix)

		// Check if the line indicates the end of the stream
		if bytes.HasPrefix(line, doneSequence) {
			break
		}

		// Create a new CompletionResponse object
		output := new(CompletionResponse)

		// Unmarshal the line into the CompletionResponse object
		if err := json.Unmarshal(line, output); err != nil {
			return fmt.Errorf("invalid json stream data: %w", err)
		}

		// Call the onData callback function with the CompletionResponse object
		onData(output)
	}

	return nil
}

//FUNCTION NOT GETTING USED
// Edits sends a request to the GPT-3 API to perform edits on a given text.
func (c *client) Edits(ctx context.Context, request EditsRequest) (*EditsResponse, error) {
	// Create a new request with the provided context, HTTP method, and request body.
	req, err := c.newRequest(ctx, "POST", "/edits", request)
	if err != nil {
		return nil, err
	}

	// Perform the request and get the response.
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	// Create a new EditsResponse object to store the response data.
	output := new(EditsResponse)

	// Parse the response and populate the output object.
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}

	return output, nil
}

//FUNCTION NOT GETTING USED
// Search sends a search request to the OpenAI API and returns the search response.
func (c *client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	// Create a new request using the provided context, HTTP method, and endpoint.
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/openai/deployments/%s/search", c.deploymentName), request)
	if err != nil {
		return nil, err
	}

	// Perform the request and get the response.
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	// Create a new SearchResponse object to hold the response data.
	output := new(SearchResponse)

	// Parse the response and populate the output object.
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}

	return output, nil
}

//FUNCTION NOT USED
// Embeddings creates text embeddings for a supplied slice of inputs with a provided model.
//
// See: https://beta.openai.com/docs/api-reference/embeddings
// It sends a POST request to the "/embeddings" endpoint with the provided request data.
// It returns the embeddings response or an error if the request fails.
func (c *client) Embeddings(ctx context.Context, request EmbeddingsRequest) (*EmbeddingsResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/embeddings", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := EmbeddingsResponse{}
	if err := getResponseObject(resp, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

//Getting called in completion, chatCompletion and multiple other functions above in this file
func (c *client) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkForSuccess(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

//Getting called in the performRequest function above
// returns an error if this response includes an error.
func checkForSuccess(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	var result APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := APIError{
			StatusCode: resp.StatusCode,
			Type:       "Unexpected",
			Message:    string(data),
		}
		return apiError
	}
	result.Error.StatusCode = resp.StatusCode
	return result.Error
}

//Getting called in this file above, in the completion and chat completion functions
func getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
}

//Getting called in the newRequest function below
// jsonBodyReader is a helper function that converts the given body interface{} into a JSON-encoded io.Reader.
func jsonBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}

	return bytes.NewBuffer(raw), nil
}

//Getting called in completion, chatCompletion and other functions above
// newRequest creates a new HTTP request with the specified method, path, and payload.
func (c *client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	// Create a JSON body reader from the payload
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}

	// Construct the request URL with the endpoint, path, and API version
	reqURL := fmt.Sprintf("%s%s?api-version=%s", c.endpoint, path, c.apiVersion)

	// Create a new HTTP request with the specified method, URL, and body reader
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set the Content-type and api-key headers
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	return req, nil
}
