package util

import (
	"bandicute-server/internal/storage/repository/post"
	"bandicute-server/internal/template"
	"bandicute-server/pkg/logger"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

const (
	referencePrefix      = "refs/heads/"
	commitFileNameFormat = "[Summary] %s - %s (%s).md"
)

type GitHubService struct {
	*github.Client
	pullRequestTemplate *template.PullRequestTemplate
}

func NewGitHubService(token string) (*GitHubService, error) {
	client := createOauth2Client(token)
	pullRequestTemplate, err := template.NewPullRequestTemplate()
	if err != nil {
		return nil, err
	}

	return &GitHubService{
		github.NewClient(client),
		pullRequestTemplate,
	}, nil
}

func (s *GitHubService) CreatePullRequestAndGetUrl(ctx context.Context, studyMemberName string, post *post.Model,
	owner, repo, filePath, content string) (string, error) {
	// 1. 깃허브 계정 핸들 가져오기
	bandicuteHandle, err := s.getBandicuteGithubHandle(ctx)
	if err != nil {
		return "", err
	}

	// 2. fork Repository
	forkedRepository, err := s.forkRepository(ctx, owner, repo, bandicuteHandle)
	if err != nil {
		return "", err
	}

	// 3. 새로운 브랜치 생성
	defaultBranchName := forkedRepository.GetDefaultBranch()
	newBranchName := createBranchName(studyMemberName)
	err = s.createNewBranchAtForkedRepository(ctx, bandicuteHandle, forkedRepository, defaultBranchName, newBranchName)
	if err != nil {
		return "", err
	}

	// 4. 브랜치에 파일 만들기
	fileNameFormat := commitFileNameFormat
	fileName := createFileName(fileNameFormat, studyMemberName, post.Title, post.PublishedAt)
	filePathAndName := fmt.Sprintf("%s/%s", filePath, fileName)
	err = s.createFileToNewBranch(ctx, bandicuteHandle, forkedRepository, newBranchName, filePathAndName, content)
	if err != nil {
		return "", err
	}

	// 5. PR Tempate 채우기
	pullRequestContent, err := s.createPullRequestContent(post, studyMemberName, content)
	if err != nil {
		return "", err
	}

	// 6. PR 생성
	pullRequestHead := fmt.Sprintf("%s:%s", bandicuteHandle, newBranchName)
	pullRequest, _, err := s.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title:               github.String(pullRequestContent.Title),
		Head:                github.String(pullRequestHead),
		Base:                github.String(defaultBranchName),
		Body:                github.String(pullRequestContent.Body),
		MaintainerCanModify: github.Bool(true),
	})

	if err != nil {
		return "", fmt.Errorf("failed to create pull request: %w", err)
	}

	logger.Info("Successfully created pull request", logger.Fields{
		"memberName": studyMemberName,
		"postTitle":  post.Title,
		"repository": owner + "/" + repo,
		"pr_url":     pullRequest.GetHTMLURL(),
	})

	return pullRequest.GetHTMLURL(), nil
}

func (s *GitHubService) getBandicuteGithubHandle(ctx context.Context) (string, error) {
	currentUser, _, err := s.Users.Get(ctx, "")
	if err != nil || currentUser == nil {
		return "", fmt.Errorf("failed to get authenticated user: %w", err)
	}

	return currentUser.GetLogin(), err
}

func (s *GitHubService) forkRepository(ctx context.Context, owner, repo, myHandle string) (*github.Repository, error) {
	forked, err := s.getForkedRepositoryIfAlreadyExists(ctx, owner, repo, myHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to check if fork exists: %w", err)
	}
	if forked != nil {
		return forked, nil
	}

	newForkedRepository, _, err := s.Repositories.CreateFork(ctx, owner, repo, &github.RepositoryCreateForkOptions{
		DefaultBranchOnly: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create fork: %w", err)
	}
	return newForkedRepository, nil
}

func (s *GitHubService) getForkedRepositoryIfAlreadyExists(ctx context.Context, owner, repo, myHandle string) (*github.Repository, error) {
	forks, _, err := s.Repositories.ListForks(ctx, owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list forks: %w", err)
	}

	for _, forked := range forks {
		if myHandle == forked.GetOwner().GetLogin() {
			return forked, nil
		}
	}
	return nil, nil
}

func (s *GitHubService) createNewBranchAtForkedRepository(ctx context.Context, myHandle string, forkedRepo *github.Repository, defaultBranchName, newBranchName string) error {
	defaultBranchReference, _, err := s.Git.GetRef(ctx, myHandle, forkedRepo.GetName(), referencePrefix+defaultBranchName)
	if err != nil {
		return fmt.Errorf("failed to get default branch reference: %w", err)
	}

	newBranchReference := &github.Reference{
		Ref:    github.String(referencePrefix + newBranchName),
		Object: defaultBranchReference.Object,
	}
	_, _, err = s.Git.CreateRef(ctx, myHandle, forkedRepo.GetName(), newBranchReference)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	return nil
}

func (s *GitHubService) createFileToNewBranch(ctx context.Context, myHandle string, forkedRepo *github.Repository, newBranchName, filePath, content string) error {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))
	_, _, err := s.Repositories.CreateFile(ctx, myHandle, forkedRepo.GetName(), filePath, &github.RepositoryContentFileOptions{
		Message: github.String("반디가 친구의 블로그에서 새로운 글을 가져왔어요!"),
		Content: []byte(encodedContent),
		Branch:  github.String(newBranchName),
	})

	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

func (s *GitHubService) createPullRequestContent(post *post.Model, studyMemberName, content string) (template.PullRequestContent, error) {
	publishedAt := post.PublishedAt.Format("2025년 01월 01일 18시 00분")
	pullRequestContent, err := s.pullRequestTemplate.FillOut(studyMemberName, post.Title, publishedAt, post.URL, content)
	if err != nil {
		return template.PullRequestContent{}, fmt.Errorf("failed to execute summaryPromptTemplate: %w", err)
	}

	return pullRequestContent, nil
}

func createFileName(fileNameFormat, memberName, postTitle string, publishedAt time.Time) string {
	formattedDate := publishedAt.Format("06-01-02")
	sanitizedTitle := sanitizeString(postTitle)
	return fmt.Sprintf(fileNameFormat, sanitizedTitle, memberName, formattedDate)
}

func createOauth2Client(token string) *http.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := oauth2.NewClient(context.Background(), tokenSource)
	return client
}

func createBranchName(memberName string) string {
	timestamp := time.Now().Format("20060102-150405")
	branchName := "summary/" + memberName + "/" + timestamp
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
