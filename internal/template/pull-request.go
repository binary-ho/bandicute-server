package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const pullRequestTemplateFileName = "pull-request-template.json"

type PullRequestTemplate struct {
	template *template.Template
}

type PullRequestContent struct {
	Title string
	Body  string
}

type pullRequestTemplateFormat struct {
	Title string `json:"title"`
	Body  struct {
		Sections []struct {
			Type    string   `json:"type"`
			Title   string   `json:"title,omitempty"`
			Content string   `json:"content,omitempty"`
			Items   []string `json:"items,omitempty"`
		} `json:"sections"`
	} `json:"body"`
}

func NewPullRequestTemplate() (*PullRequestTemplate, error) {
	instance, err := getPullRequestTemplateInstance()
	if err != nil {
		return nil, err
	}
	return &PullRequestTemplate{template: instance}, nil
}

func (t *PullRequestTemplate) FillOut(memberName, postTitle, publishedAt, postUrl, summary string) (PullRequestContent, error) {
	var templateBuf bytes.Buffer
	err := t.template.Execute(&templateBuf, map[string]interface{}{
		"member_name":  memberName,
		"post_title":   postTitle,
		"published_at": publishedAt,
		"post_url":     postUrl,
		"summary":      summary,
	})

	if err != nil {
		return PullRequestContent{}, fmt.Errorf("failed to execute PullRequestTemplate: %w", err)
	}

	parts := strings.SplitN(templateBuf.String(), "\n---\n", 2)
	if len(parts) != 2 {
		return PullRequestContent{}, fmt.Errorf("invalid summaryPromptTemplate output format")
	}

	return PullRequestContent{
		Title: strings.TrimSpace(parts[0]),
		Body:  strings.TrimSpace(parts[1]),
	}, nil
}

var pullRequestTemplateInstance *template.Template

func getPullRequestTemplateInstance() (*template.Template, error) {
	if pullRequestTemplateInstance != nil {
		return pullRequestTemplateInstance, nil
	}

	var format pullRequestTemplateFormat
	err := parseTemplateByFormat(format, pullRequestTemplateFileName)
	result := buildTemplate(format)

	pullRequestTemplateInstance, err = template.New(pullRequestTemplateFileName).Parse(result)
	if err != nil {
		return nil, fmt.Errorf("failed to create pullRequestTemplateInstance : %w", err)
	}
	return pullRequestTemplateInstance, nil
}

func buildTemplate(format pullRequestTemplateFormat) string {
	pullRequestContent := buildContent(format)
	return format.Title + "\n---\n" + pullRequestContent
}

func buildContent(format pullRequestTemplateFormat) string {
	var contentBuilder strings.Builder
	for _, section := range format.Body.Sections {
		switch section.Type {
		case "header", "footer":
			contentBuilder.WriteString(section.Content + "\n\n")
		case "info", "summary", "recommendation":
			if section.Title != "" {
				contentBuilder.WriteString("## " + section.Title + "\n\n")
			}
			if len(section.Items) > 0 {
				for _, item := range section.Items {
					contentBuilder.WriteString("- " + item + "\n")
				}
				contentBuilder.WriteString("\n")
			}
			if section.Content != "" {
				contentBuilder.WriteString(section.Content + "\n\n")
			}
		}
	}
	return contentBuilder.String()
}
