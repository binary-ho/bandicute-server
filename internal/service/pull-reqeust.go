package service

import (
	"bandicute-server/internal/service/request"
	pullRequest "bandicute-server/internal/storage/repository/pull-request"
	"bandicute-server/internal/util"
	"bandicute-server/pkg/logger"
	"fmt"
	"strings"
	"time"
)

type PullRequestOpener struct {
	githubService         *util.GitHubService
	pullRequestRepository pullRequest.Repository
}

func NewPullRequestOpener(
	githubService *util.GitHubService,
	pullRequestRepository pullRequest.Repository) *PullRequestOpener {
	return &PullRequestOpener{
		githubService:         githubService,
		pullRequestRepository: pullRequestRepository,
	}
}

func (p *PullRequestOpener) OpenPullRequest(req request.OpenPullRequest) {
	post := req.Post
	owner, repositoryName, err := parseOwnerAndRepo(req.Repository)
	if err != nil {
		logger.Error("Failed to parse owner and repository", logger.Fields{
			"repo":  req.Repository,
			"error": err.Error(),
		})
		return
	}

	pullRequestRul, err := p.githubService.CreatePullRequestAndGetUrl(
		req.Context, req.MemberName, post,
		owner, repositoryName, req.FilePath, req.Summary,
	)

	if err != nil {
		p.pullRequestRepository.Create(req.Context, getEmptyPullRequest(post.ID, req.StudyId))
		logger.Error("Failed to create PR", logger.Fields{
			"post":   post,
			"repo":   req.Repository,
			"member": req.MemberName,
			"error":  err.Error(),
		})
		return
	}

	_, err = p.pullRequestRepository.Create(req.Context, &pullRequest.Model{
		BlogPostID: post.ID,
		StudyID:    req.StudyId,
		PrUrl:      pullRequestRul,
		IsOpened:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	if err != nil {
		logger.Error("Failed to save PR to DB", logger.Fields{
			"post":   post,
			"repo":   req.Repository,
			"member": req.MemberName,
			"error":  err.Error(),
		})
		return
	}
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

func getEmptyPullRequest(postId, studyId string) *pullRequest.Model {
	return &pullRequest.Model{
		BlogPostID: postId,
		StudyID:    studyId,
		PrUrl:      "",
		IsOpened:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
