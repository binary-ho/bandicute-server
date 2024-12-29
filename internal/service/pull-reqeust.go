package service

import (
	"bandicute-server/internal/service/request"
	pullRequest "bandicute-server/internal/storage/repository/pull-request"
	"bandicute-server/internal/util"
	"bandicute-server/pkg/logger"
	"time"
)

type PullRequestOpener struct {
	pullRequestRepository pullRequest.Repository
	githubService         util.GitHubService
}

func (p *PullRequestOpener) OpenPullRequest(req request.OpenPullRequest) {
	post := req.Post
	pullRequestRul, err := p.githubService.CreatePR(req.Context, post, req.Repository, req.MemberName, req.Summary)
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
