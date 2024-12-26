package service

import (
	"bandicute-server/internal/storage/repository/member"
	"bandicute-server/internal/storage/repository/post"
	"bandicute-server/internal/storage/repository/study"
	"bandicute-server/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type GitHubPRService struct {
	client     *github.Client
	prTemplate *template.Template
}

func NewGitHubPRService(token string) (*GitHubPRService, error) {
	// Load and parse PR template
	templateBytes, err := os.ReadFile("internal/templates/pr-template.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read PR template: %w", err)
	}

	var prTemplate struct {
		Title string `json:"title"`
		Body  struct {
			Sections []struct {
				Type    string   `json:"type"`
				Title   string   `json:"title,omitempty"`
				Content string   `json:"content,omitempty"`
				Items   []string `json:"items,omitempty"`
			} `json:"sections"`
		} `json:"body"`
		Variables []string `json:"variables"`
	}
	if err := json.Unmarshal(templateBytes, &prTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse PR template: %w", err)
	}

	// Create PR body template
	var bodyBuilder strings.Builder
	for _, section := range prTemplate.Body.Sections {
		switch section.Type {
		case "header", "footer":
			bodyBuilder.WriteString(section.Content + "\n\n")
		case "info", "summary", "recommendation":
			if section.Title != "" {
				bodyBuilder.WriteString("## " + section.Title + "\n\n")
			}
			if len(section.Items) > 0 {
				for _, item := range section.Items {
					bodyBuilder.WriteString("- " + item + "\n")
				}
				bodyBuilder.WriteString("\n")
			}
			if section.Content != "" {
				bodyBuilder.WriteString(section.Content + "\n\n")
			}
		}
	}

	tmpl, err := template.New("pr").Parse(prTemplate.Title + "\n---\n" + bodyBuilder.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &GitHubPRService{
		client:     github.NewClient(tc),
		prTemplate: tmpl,
	}, nil
}

func (s *GitHubPRService) CreatePR(ctx context.Context, study *study.Model, member *member.Model, post *post.Model, summary string) (string, error) {
	// Parse storage Repository
	repoURL := study.Repository
	parts := strings.Split(repoURL, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid storage Repository: %s", repoURL)
	}
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Get default branch
	repository, _, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get storage: %w", err)
	}
	defaultBranch := repository.GetDefaultBranch()

	// Create branch name
	timestamp := time.Now().Format("20060102-150405")
	sanitizedTitle := sanitizeString(post.Title)
	branchName := fmt.Sprintf("post/%s-%s", timestamp, sanitizedTitle)

	// Get the default branch reference
	ref, _, err := s.client.Git.GetRef(ctx, owner, repo, "refs/heads/"+defaultBranch)
	if err != nil {
		return "", fmt.Errorf("failed to get reference: %w", err)
	}

	// Create a new branch
	newRef := &github.Reference{
		Ref:    github.String("refs/heads/" + branchName),
		Object: ref.Object,
	}
	_, _, err = s.client.Git.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	// Format date
	publishedAt := post.PublishedAt.Format("2006년 01월 02일")

	// Execute template
	var bodyBuf bytes.Buffer
	err = s.prTemplate.Execute(&bodyBuf, map[string]interface{}{
		"member_name":  member.Name,
		"post_title":   post.Title,
		"published_at": publishedAt,
		"post_url":     post.URL,
		"summary":      summary,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Split title and body
	parts = strings.SplitN(bodyBuf.String(), "\n---\n", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid template output format")
	}
	title := strings.TrimSpace(parts[0])
	body := strings.TrimSpace(parts[1])

	// Create pull request
	pr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(branchName),
		Base:  github.String(defaultBranch),
		Body:  github.String(body),
	}

	pullRequest, _, err := s.client.PullRequests.Create(ctx, owner, repo, pr)
	if err != nil {
		return "", fmt.Errorf("failed to create pull request: %w", err)
	}

	logger.Info("Successfully created pull request", logger.Fields{
		"pr_number": pullRequest.GetNumber(),
		"pr_url":    pullRequest.GetHTMLURL(),
	})

	return pullRequest.GetHTMLURL(), nil
}

func sanitizeString(s string) string {
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return strings.ToLower(s)
}
