package template

import (
	"bytes"
	"fmt"
	"text/template"
)

const summaryPromptTemplateFileName = "summary-prompt.json"

type SummaryPromptTemplate struct {
	template *template.Template
}

type summaryPromptTemplateFormat struct {
	Template string `json:"template"`
}

func NewSummaryPromptTemplate() (*SummaryPromptTemplate, error) {
	instance, err := getSummaryPromptTemplateInstance()
	if err != nil {
		return nil, err
	}
	return &SummaryPromptTemplate{template: instance}, nil
}

func (t *SummaryPromptTemplate) FillOut(title, content string) (string, error) {
	var templateBuf bytes.Buffer
	err := t.template.Execute(&templateBuf, map[string]interface{}{
		"title":   title,
		"content": content,
	})

	if err != nil {
		return "", fmt.Errorf("failed to fill out summaryPromptTemplate: %w", err)
	}
	return templateBuf.String(), nil
}

var summaryPromptTemplateInstance *template.Template

func getSummaryPromptTemplateInstance() (*template.Template, error) {
	if summaryPromptTemplateInstance != nil {
		return summaryPromptTemplateInstance, nil
	}

	format := &summaryPromptTemplateFormat{}
	err := parseTemplateByFormat(format, summaryPromptTemplateFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to parseFeed summary prompt promptTemplate : %w", err)
	}

	summaryPromptTemplateInstance, err = template.New(summaryPromptTemplateFileName).Parse(format.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to create summaryPromptTemplate : %w", err)
	}
	return summaryPromptTemplateInstance, nil
}
