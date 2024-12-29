package service

import (
	"bandicute-server/internal/service/channel"
	"bandicute-server/internal/service/request"
	"bandicute-server/internal/storage/repository/member"
	"bandicute-server/internal/storage/repository/post"
	"bandicute-server/internal/storage/repository/summary"
	"bandicute-server/internal/util"
	"bandicute-server/pkg/logger"
	"context"
	"fmt"
)

type Parser struct {
	parser            *util.PostParser
	memberRepository  member.Repository
	postRepository    post.Repository
	summaryRepository summary.Repository
}

func NewParser(
	parser *util.PostParser,
	memberRepository member.Repository,
	postRepository post.Repository,
	summaryRepository summary.Repository,
) *Parser {
	return &Parser{
		parser:            parser,
		memberRepository:  memberRepository,
		postRepository:    postRepository,
		summaryRepository: summaryRepository,
	}
}

func (w *Parser) ParseRecentPostByMember(ctx context.Context, memberId string, summarizeRequestChannel *channel.SummarizeRequest) {
	// 1. Get StudyMember
	member, err := w.memberRepository.GetById(ctx, memberId)
	if err != nil {
		logger.Error("Failed to get study member", logger.Fields{
			"memberId": memberId,
			"error":    err.Error(),
		})
		return
	}

	// 2. Parse Member's Blog
	recentPosts, err := w.parseRecentPostsByMember(ctx, err, member)
	if err != nil {
		logger.Error("Failed to parse recent posts", logger.Fields{
			"member": member,
			"error":  err.Error(),
		})
		return
	}

	for _, eachPost := range recentPosts {
		go w.createPostAndRequestSummarize(ctx, eachPost, summarizeRequestChannel)
	}
	return
}

func (w *Parser) parseRecentPostsByMember(ctx context.Context, err error, member *member.Model) ([]*post.Model, error) {
	// 1. Parse post
	posts, err := w.parser.Parse(ctx, member.Blog)
	if err != nil {
		return nil, fmt.Errorf("failed to parseFeed post: %w", err)
	}

	// 2. Filter recent post
	latestPost, err := w.postRepository.GetLatestByMemberId(ctx, member.ID)
	return util.FilterRecentPost(latestPost, posts), nil
}

func (w *Parser) createPostAndRequestSummarize(ctx context.Context, post *post.Model, requestChannel *channel.SummarizeRequest) {
	// 1. Save post
	savedPost, err := w.postRepository.Create(ctx, post)
	if err != nil {
		logger.Error("Failed to create post", logger.Fields{
			"post":  post,
			"error": err.Error(),
		})
		return
	}

	// TODO: 2. Save empty summary
	// + 포스트 저장과, Summarize 생성은 한 트랜잭션에 묶여야 한다. 또한 현재 실패에 대한 고려가 없음

	// 3. Request summarize
	*requestChannel <- request.Summarize{
		Context: ctx,
		Post:    savedPost,
	}
}
