package service

import (
	"bandicute-server/internal/service/channel"
	"bandicute-server/internal/service/request"
	"bandicute-server/internal/storage/repository/member"
	"bandicute-server/internal/storage/repository/post"
	pullRequest "bandicute-server/internal/storage/repository/pull-request"
	"bandicute-server/internal/storage/repository/study"
	"bandicute-server/internal/storage/repository/summary"
	"bandicute-server/internal/util"
	"bandicute-server/pkg/logger"
	"time"
)

type Summarizer struct {
	summarizer            *util.Summarizer
	memberRepository      member.Repository
	summaryRepository     summary.Repository
	studyRepository       study.Repository
	pullRequestRepository pullRequest.Repository
}

func (s *Summarizer) Summarize(req request.Summarize, openPullRequestRequestChannel *channel.OpenPullRequestRequest) {
	post := req.Post

	// 1. 요약 요청
	summaryContent, err := s.summarizer.Summarize(req.Context, post.Title, post.Content)
	if err != nil {
		_, err = s.summaryRepository.Create(req.Context, getEmptySummary(post))
		logger.Error("Failed to summarize post", logger.Fields{
			"post":  post,
			"error": err.Error(),
		})
		return
	}

	// 2. 요약문 저장
	// TODO: Create가 아니라, 이미 있는걸 Update할 수 있어야 한다. 일종의 Merge가 구현되어야 한다.
	_, err = s.summaryRepository.Create(req.Context, &summary.Model{
		BlogPostID:   post.ID,
		Summary:      summaryContent,
		IsSummarized: true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		logger.Error("Failed to create summary", logger.Fields{
			"post":  post,
			"error": err.Error(),
		})
		return
	}

	// TODO: 3. member가 가입되어 있는 모든 study에 대해 Empty PR을 생성
	// 4. Request Open Pull Request
	s.requestOpenPullRequest(req, post, summaryContent, openPullRequestRequestChannel)
}

func (s *Summarizer) requestOpenPullRequest(req request.Summarize, post *post.Model, summaryContent string, openPullRequestRequestChannel *channel.OpenPullRequestRequest) {
	studies, err := s.studyRepository.GetAllByMemberId(req.Context, post.MemberID)
	if err != nil {
		logger.Error("Failed to get studies in 'Request Open Pull Request'", logger.Fields{
			"memberId": post.MemberID,
			"error":    err.Error(),
		})
		return
	}

	member, err := s.memberRepository.GetById(req.Context, post.MemberID)
	if err != nil {
		logger.Error("Failed to get member in 'Request Open Pull Request'", logger.Fields{
			"memberId": post.MemberID,
			"error":    err.Error(),
		})
		return
	}

	for _, study := range studies {
		*openPullRequestRequestChannel <- request.OpenPullRequest{
			Context:    req.Context,
			Post:       post,
			Repository: study.Repository,
			MemberName: member.Name,
			Summary:    summaryContent,
			StudyId:    study.ID,
		}
	}
}

func getEmptySummary(post *post.Model) *summary.Model {
	return &summary.Model{
		BlogPostID:   post.ID,
		Summary:      "",
		IsSummarized: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
