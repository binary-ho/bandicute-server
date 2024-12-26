package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

type Summarizer struct {
	client *openai.Client
	prompt *template.Template
}

func NewGPTSummarizer(apiKey string) (*Summarizer, error) {
	// Load and parse summary prompt template
	promptBytes, err := os.ReadFile("internal/templates/summary-prompt.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read summary prompt template: %w", err)
	}

	var promptTemplate struct {
		Template  string   `json:"template"`
		Variables []string `json:"variables"`
	}
	if err := json.Unmarshal(promptBytes, &promptTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse summary prompt template: %w", err)
	}

	tmpl, err := template.New("summary").Parse(promptTemplate.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Summarizer{
		client: openai.NewClient(apiKey),
		prompt: tmpl,
	}, nil
}

func (s *Summarizer) Summarize(ctx context.Context, title, content string) (string, error) {
	var promptBuf bytes.Buffer
	err := s.prompt.Execute(&promptBuf, map[string]interface{}{
		"title":   title,
		"content": content,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: promptBuf.String(),
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
