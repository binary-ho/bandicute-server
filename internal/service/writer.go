package service

import (
	"bandicute-server/internal/storage/repository/member"
	"bandicute-server/internal/storage/repository/post"
	pullRequest "bandicute-server/internal/storage/repository/pull-request"
	"bandicute-server/internal/storage/repository/study"
	studyMember "bandicute-server/internal/storage/repository/study-member"
	"bandicute-server/internal/storage/repository/summary"
	"bandicute-server/pkg/logger"
	"context"
	"fmt"
	"strings"
	"time"
)

type RecentPostSummaryWriter struct {
	parser                Parser
	summarizer            *Summarizer
	prService             *GitHubPRService
	memberRepository      member.Repository
	studyRepository       study.Repository
	studyMemberRepository studyMember.Repository
	postRepository        post.Repository
	summaryRepository     summary.Repository
	pullRequestRepository pullRequest.Repository
}

func NewRecentPostSummaryWriter(parser Parser, summarizer *Summarizer, prService *GitHubPRService) (*RecentPostSummaryWriter, error) {
	return &RecentPostSummaryWriter{
		parser:     parser,
		summarizer: summarizer,
		prService:  prService,
	}, nil
}

func (w *RecentPostSummaryWriter) ProcessStudy(ctx context.Context, studyID string) error {
	// Get study
	study, err := w.studyRepository.GetStudy(ctx, studyID)
	if err != nil {
		return fmt.Errorf("failed to get study: %w", err)
	}

	// Get study members
	members, err := w.studyMemberRepository.GetStudyMembers(ctx, studyID)
	if err != nil {
		return fmt.Errorf("failed to get study members: %w", err)
	}

	logger.Info("Processing study", logger.Fields{
		"study_id":      studyID,
		"member_count":  len(members),
		"github_repo":   study.Repository,
		"github_branch": study.Branch,
		"github_dir":    study.Directory,
	})

	// Process each member
	for _, member := range members {
		if err := w.processMember(ctx, study, member); err != nil {
			logger.Error("Failed to process member", logger.Fields{
				"study_id":  studyID,
				"member_id": member.MemberID,
				"error":     err.Error(),
			})
			// Continue with other members even if one fails
			continue
		}
	}

	return nil
}

func (w *RecentPostSummaryWriter) ProcessStudyMember(ctx context.Context, studyMemberID string) error {
	// Get study studyMember
	studyMember, err := w.GetStudyMember(ctx, studyMemberID)
	if err != nil {
		return fmt.Errorf("failed to get study studyMember: %w", err)
	}

	// Get study
	study, err := w.studyRepository.GetStudy(ctx, studyMember.StudyId)
	if err != nil {
		return fmt.Errorf("failed to get study: %w", err)
	}

	return w.processMember(ctx, study, studyMember)
}

func (w *RecentPostSummaryWriter) GetStudyMember(ctx context.Context, studyMemberID string) (*studyMember.Model, error) {
	studyMember, err := w.studyMemberRepository.GetStudyMember(ctx, studyMemberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study studyMember: %w", err)
	}
	return studyMember, nil
}

func (w *RecentPostSummaryWriter) processMember(ctx context.Context, study *study.Model, studyMember *studyMember.Model) error {
	// Get member
	member, err := w.memberRepository.GetMember(ctx, studyMember.MemberID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	logger.Info("Processing member", logger.Fields{
		"member_id":    member.ID,
		"tistory_blog": member.Blog,
		"name":         member.Name,
	})

	// Parse post
	posts, err := w.parser.ParseBlog(ctx, member.Blog)
	if err != nil {
		return fmt.Errorf("failed to parse post: %w", err)
	}

	// Process each post
	for _, post := range posts {
		if err := w.processPost(ctx, study, member, post, posts); err != nil {
			logger.Error("Failed to process post", logger.Fields{
				"member_id": member.ID,
				"post_url":  post.URL,
				"error":     err.Error(),
			})
			// Continue with other posts even if one fails
			continue
		}
	}

	return nil
}

func (w *RecentPostSummaryWriter) processPost(ctx context.Context, study *study.Model, member *member.Model, post *post.Model, posts []*post.Model) error {
	// 1. 이미 처리된 포스트인지 GUID로 확인
	existingPost, err := w.postRepository.GetPostByGUID(ctx, post.GUID)
	if err == nil {
		logger.Info("Post already exists", logger.Fields{
			"post_id": existingPost.ID,
			"url":     existingPost.URL,
		})
		return nil
	}

	// 2. 멤버의 최신 포스트 확인
	latestPost, err := w.postRepository.GetLatestPost(ctx, member.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "no post posts found") {
			return fmt.Errorf("failed to get latest post post: %w", err)
		}
		// 첫 포스트인 경우, RSS 피드의 가장 최신 글만 처리
		logger.Info("No existing posts found for member, processing only the latest post", logger.Fields{
			"member_id": member.ID,
			"post_url":  post.URL,
		})
		if !isLatestInFeed(post, posts) {
			logger.Info("Skipping non-latest post for new member", logger.Fields{
				"post_url": post.URL,
			})
			return nil
		}
	} else {
		// 기존 포스트가 있는 경우, 발행 시간 비교
		if !post.PublishedAt.After(latestPost.PublishedAt) {
			logger.Info("Skipping old post", logger.Fields{
				"post_url":                 post.URL,
				"published_at":             post.PublishedAt,
				"latest_post_published_at": latestPost.PublishedAt,
			})
			return nil
		}
	}

	// Set member ID
	post.MemberID = member.ID

	// Create post
	if err := w.postRepository.CreatePost(ctx, post); err != nil {
		return fmt.Errorf("failed to create post post: %w", err)
	}

	// Generate summary using GPT
	postContent, err := w.summarizer.Summarize(ctx, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	// Create post summary
	summary := &summary.Model{
		BlogPostID:   post.ID,
		Summary:      postContent,
		IsSummarized: true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := w.summaryRepository.CreatePostSummary(ctx, summary); err != nil {
		return fmt.Errorf("failed to create post summary: %w", err)
	}

	// Create GitHub PR
	prURL, err := w.prService.CreatePR(ctx, study, member, post, summary.Summary)
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	// Create pull request record
	pr := &pullRequest.Model{
		BlogPostID: post.ID,
		StudyID:    study.ID,
		PrUrl:      prURL,
		IsOpened:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := w.pullRequestRepository.CreatePullRequest(ctx, pr); err != nil {
		return fmt.Errorf("failed to create pull request record: %w", err)
	}

	logger.Info("Successfully processed post", logger.Fields{
		"post_url": post.URL,
		"pr_url":   prURL,
	})

	return nil
}

// isLatestInFeed checks if the given post is the latest in the feed
func isLatestInFeed(post *post.Model, posts []*post.Model) bool {
	for _, p := range posts {
		if p.PublishedAt.After(post.PublishedAt) {
			return false
		}
	}
	return true
}
