package cli

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/PullRequestInc/go-gpt3"
	azureopenai "github.com/akhilsharma90/terraform-assistant/pkg/gpt3"
	"github.com/akhilsharma90/terraform-assistant/pkg/utils"
	"github.com/pkg/errors"
)

var errResp = errors.New("invalid response")

// openaiGptCompletion generates a GPT-3 completion using the OpenAI API.
func (c *oaiClients) openaiGptCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	// Make a completion request to the OpenAI API
	resp, err := c.openAIClient.CompletionWithEngine(ctx, *openAIDeploymentName, openai.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           utils.ToPtr(1),
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error openai completion: %w", err)
	}

	// Check if the response contains exactly one choice
	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	// Return the generated completion text
	return resp.Choices[0].Text, nil
}

func (c *oaiClients) openaiGptChatCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.openAIClient.ChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: *openAIDeploymentName,
		Messages: []openai.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           1,
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error openai gpt completion: %w", err)
	}

	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *oaiClients) azureGptCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.azureClient.Completion(ctx, azureopenai.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           utils.ToPtr(1),
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error azure completion: %w", err)
	}

	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

func (c *oaiClients) azureGptChatCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.azureClient.ChatCompletion(ctx, azureopenai.ChatCompletionRequest{
		Model: *openAIDeploymentName,
		Messages: []azureopenai.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           1,
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error azure chatgpt completion: %w", err)
	}

	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil
}
