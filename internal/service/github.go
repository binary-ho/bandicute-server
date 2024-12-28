package service

import (
	"bandicute-server/internal/storage/repository/member"
	"bandicute-server/internal/storage/repository/post"
	"bandicute-server/internal/storage/repository/study"
	"bandicute-server/internal/template"
	"bandicute-server/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

const ReferencePrefix = "refs/heads/"

type GitHubPRService struct {
	client              *github.Client
	pullRequestTemplate *template.PullRequestTemplate
}

func NewGitHubPRService(token string) (*GitHubPRService, error) {
	client := createOauth2Client(token)
	pullRequestTemplate, err := template.NewPullRequestTemplate()
	if err != nil {
		return nil, err
	}

	return &GitHubPRService{
		client:              github.NewClient(client),
		pullRequestTemplate: pullRequestTemplate,
	}, nil
}

func (s *GitHubPRService) CreatePR(ctx context.Context, study *study.Model, member *member.Model, post *post.Model, summary string) (string, error) {
	// 1. Parse storage Repository
	owner, repo, err := parseOwnerAndRepo(study.Repository)
	if err != nil {
		return "", err
	}

	// 2. Get default Branch And Reference
	repository, _, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get storage: %w", err)
	}
	defaultBranch := repository.GetDefaultBranch()
	defaultBranchReference, _, err := s.client.Git.GetRef(ctx, owner, repo, ReferencePrefix+defaultBranch)
	if err != nil {
		return "", fmt.Errorf("failed to get reference: %w", err)
	}

	// 3. Create a new Branch And Reference
	newBranchName := createBranchName(post.Title)
	newBranchReference := &github.Reference{
		Ref:    github.String(ReferencePrefix + newBranchName),
		Object: defaultBranchReference.Object,
	}

	_, _, err = s.client.Git.CreateRef(ctx, owner, repo, newBranchReference)
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	// 4. PR Tempate 채우기
	publishedAt := post.PublishedAt.Format("2006년 01월 02일")
	pullRequestContent, err := s.pullRequestTemplate.FillOut(member.Name, post.Title, publishedAt, post.URL, summary)
	if err != nil {
		return "", fmt.Errorf("failed to execute summaryPromptTemplate: %w", err)
	}

	// 5. PR 생성하기
	pullRequest, _, err := s.client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: github.String(pullRequestContent.Title),
		Head:  github.String(newBranchName),
		Base:  github.String(defaultBranch),
		Body:  github.String(pullRequestContent.Body),
	})

	if err != nil {
		return "", fmt.Errorf("failed to create pull request: %w", err)
	}

	logger.Info("Successfully created pull request", logger.Fields{
		"pr_number": pullRequest.GetNumber(),
		"pr_url":    pullRequest.GetHTMLURL(),
	})

	return pullRequest.GetHTMLURL(), nil
}

func createOauth2Client(token string) *http.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := oauth2.NewClient(context.Background(), tokenSource)
	return client
}

func parseOwnerAndRepo(repository string) (owner, repo string, err error) {
	parts := strings.Split(repository, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid storage Repository: %s", repository)
	}

	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]
	return owner, repo, nil
}

func createBranchName(postTitle string) string {
	timestamp := time.Now().Format("20060102-150405")
	sanitizedTitle := sanitizeString(postTitle)
	branchName := fmt.Sprintf("post/%s-%s", timestamp, sanitizedTitle)
	return branchName
}

func sanitizeString(s string) string {
	// 1. 공백 지워버려
	s = strings.ReplaceAll(s, " ", "_")

	// 2. 특수문자 지우기
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return strings.ToLower(s)
}
