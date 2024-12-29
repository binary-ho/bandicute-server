package util

import (
	"bandicute-server/internal/template"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type PostSummarizer struct {
	openaiClient          *openai.Client
	summaryPromptTemplate *template.SummaryPromptTemplate
}

func NewPostSummarizer(apiKey string) (*PostSummarizer, error) {
	summaryPromptTemplate, err := template.NewSummaryPromptTemplate()
	if err != nil {
		return nil, err
	}

	return &PostSummarizer{
		openaiClient:          openai.NewClient(apiKey),
		summaryPromptTemplate: summaryPromptTemplate,
	}, nil
}

func (s *PostSummarizer) Summarize(ctx context.Context, title, content string) (string, error) {
	summaryPrompt, err := s.summaryPromptTemplate.FillOut(title, content)
	if err != nil {
		return "", err
	}

	request := getSummaryRequest(summaryPrompt)
	response, err := s.openaiClient.CreateChatCompletion(ctx, request)
	if err = validateResponse(response, err); err != nil {
		return "", err
	}

	return extractContent(response), nil
}

func getSummaryRequest(summaryPrompt string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleUser,
			Content: summaryPrompt,
		}},
	}
}

func validateResponse(response openai.ChatCompletionResponse, err error) error {
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(response.Choices) == 0 {
		return fmt.Errorf("no completion choices returned")
	}
	return nil
}

func extractContent(response openai.ChatCompletionResponse) string {
	return response.Choices[0].Message.Content
}
