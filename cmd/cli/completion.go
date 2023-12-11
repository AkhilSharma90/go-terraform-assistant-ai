package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	openai "github.com/PullRequestInc/go-gpt3"
	azureopenai "github.com/akhilsharma90/terraform-assistant/pkg/gpt3"
	"github.com/pkg/errors"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
)

// Constant for user role
const userRole = "user"

var (
	// Map to hold the maximum tokens allowed for different GPT models
	maxTokensMap = map[string]int{
		"code-davinci-002":   8001,
		"text-davinci-003":   4097,
		"gpt-3.5-turbo-0301": 4096,
		"gpt-3.5-turbo":      4096,
		"gpt-35-turbo-0301":  4096, // for azure
		"gpt-4-0314":         8192,
		"gpt-4-32k-0314":     8192,
	}
	// Error for invalid max tokens
	errToken = errors.New("invalid max tokens")
)

// Struct to hold the clients for Azure and OpenAI
type oaiClients struct {
	azureClient  azureopenai.Client
	openAIClient openai.Client
}

// Function to create new OpenAI and Azure clients
func newOAIClients() (oaiClients, error) {
	var (
		oaiClient   openai.Client
		azureClient azureopenai.Client
		err         error
	)

	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		// Create a new OpenAI client
		oaiClient = openai.NewClient(*openAIAPIKey)
	} else {
		// Validate the deployment name
		re := regexp.MustCompile(`^[a-zA-Z0-9]+([_-]?[a-zA-Z0-9]+)*$`)
		if !re.MatchString(*openAIDeploymentName) {
			return oaiClients{}, errors.New("azure openai deployment can only include alphanumeric characters, '_,-', and can't end with '_' or '-'")
		}

		// Create a new Azure client
		azureClient, err = azureopenai.NewClient(*azureOpenAIEndpoint, *openAIAPIKey, *openAIDeploymentName)
		if err != nil {
			return oaiClients{}, fmt.Errorf("error create Azure client: %w", err)
		}
	}

	// Create a new oaiClients struct with the created clients
	clients := oaiClients{
		azureClient:  azureClient,
		openAIClient: oaiClient,
	}

	return clients, nil
}

// completion is a function that generates completions for a given prompt and deployment configuration.
// It uses the provided OpenAI and Azure clients to make API calls for completion generation.
func completion(ctx context.Context, client oaiClients, prompts []string, deploymentName string, subcommand string) (string, error) {
	// Set the temperature for completion generation
	temp := float32(*temperature)

	// Calculate the maximum tokens allowed for the given deployment name
	maxTokens, err := calculateMaxTokens(prompts, deploymentName)
	if err != nil {
		return "", fmt.Errorf("error calculate max token: %w", err)
	}

	// Build the prompt string
	var prompt strings.Builder
	_, err = fmt.Fprint(&prompt, subcommand)
	if err != nil {
		return "", fmt.Errorf("error prompt string builder: %w", err)
	}

	// Append each prompt to the prompt string
	for _, p := range prompts {
		_, err = fmt.Fprintf(&prompt, "%s\n", p)
		if err != nil {
			return "", fmt.Errorf("error range prompt: %w", err)
		}
	}

	// Check if Azure OpenAI endpoint is not set
	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		// Check if the deployment name is GPT Turbo or GPT-4
		if isGptTurbo(deploymentName) || isGpt4(deploymentName) {
			// Generate completion using OpenAI GptChat completion API
			resp, err := client.openaiGptChatCompletion(ctx, prompt, maxTokens, temp)
			if err != nil {
				return "", fmt.Errorf("error openai GptChat completion: %w", err)
			}

			return resp, nil
		}

		// Generate completion using OpenAI Gpt completion API
		resp, err := client.openaiGptCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", fmt.Errorf("error openai Gpt completion: %w", err)
		}

		return resp, nil
	}

	// Check if the deployment name is GPT Turbo 3.5 or GPT-4
	if isGptTurbo35(deploymentName) || isGpt4(deploymentName) {
		// Generate completion using Azure GptChat completion API
		resp, err := client.azureGptChatCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", fmt.Errorf("error azure GptChat completion: %w", err)
		}

		return resp, nil
	}

	// Generate completion using Azure Gpt completion API
	resp, err := client.azureGptCompletion(ctx, prompt, maxTokens, temp)
	if err != nil {
		return "", fmt.Errorf("error azure Gpt completion: %w", err)
	}

	return resp, nil
}

// isGptTurbo is a function that checks deployment names.
func isGptTurbo(deploymentName string) bool {
	return deploymentName == "gpt-3.5-turbo-0301" || deploymentName == "gpt-3.5-turbo"
}

// isGptTurbo35 is a function that checks deployment names.
func isGptTurbo35(deploymentName string) bool {
	return deploymentName == "gpt-35-turbo-0301" || deploymentName == "gpt-35-turbo"
}

// isGpt4 is a function that checks deployment names.
func isGpt4(deploymentName string) bool {
	return deploymentName == "gpt-4-0314" || deploymentName == "gpt-4-32k-0314"
}

// calculateMaxTokens is a function that calculates the maximum tokens allowed for a given deployment name.
func calculateMaxTokens(prompts []string, deploymentName string) (*int, error) {
	// Get the maximum tokens allowed for the deploymentName from the maxTokensMap
	maxTokensFinal, ok := maxTokensMap[deploymentName]
	if !ok {
		return nil, errors.Wrapf(errToken, "deploymentName %q not found in max tokens map", deploymentName)
	}

	// If a custom maxTokens value is provided, override the value from the map
	if *maxTokens > 0 {
		maxTokensFinal = *maxTokens
	}

	// Create a new gptEncoder
	encoder, err := gptEncoder.NewEncoder()
	if err != nil {
		return nil, fmt.Errorf("error encode gpt: %w", err)
	}

	// Start at 100 since the encoder at times doesn't get it exactly correct
	totalTokens := 100

	// Encode each prompt and calculate the total number of tokens
	for _, prompt := range prompts {
		tokens, err := encoder.Encode(prompt)
		if err != nil {
			return nil, fmt.Errorf("error encode prompt: %w", err)
		}

		totalTokens += len(tokens)
	}

	// Calculate the remaining tokens by subtracting the total tokens from the maximum tokens allowed
	remainingTokens := maxTokensFinal - totalTokens

	return &remainingTokens, nil
}
