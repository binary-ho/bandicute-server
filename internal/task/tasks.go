package task

import (
	"bandicute-server/internal/service"
	"bandicute-server/internal/service/channel"
)

type Handler struct {
	parser                            *service.Parser
	summarizer                        *service.Summarizer
	pullRequestOpener                 *service.PullRequestOpener
	parsePostByMemberIdRequestChannel *channel.ParsePostByMemberIdRequest
	summarizeRequestChannel           *channel.SummarizeRequest
	openPullRequestRequestChannel     *channel.OpenPullRequestRequest
}

func NewHandler(
	parser *service.Parser,
	summarizer *service.Summarizer,
	pullRequestOpener *service.PullRequestOpener,
	parsePostByMemberIdRequestChannel *channel.ParsePostByMemberIdRequest,
	summarizeRequestChannel *channel.SummarizeRequest,
	openPullRequestRequestChannel *channel.OpenPullRequestRequest,
) *Handler {
	return &Handler{
		parser:                            parser,
		summarizer:                        summarizer,
		pullRequestOpener:                 pullRequestOpener,
		parsePostByMemberIdRequestChannel: parsePostByMemberIdRequestChannel,
		summarizeRequestChannel:           summarizeRequestChannel,
		openPullRequestRequestChannel:     openPullRequestRequestChannel,
	}
}

func (h *Handler) Run() {
	go func() {
		for {
			select {
			case request := <-*h.parsePostByMemberIdRequestChannel:
				go h.parser.ParseRecentPostByMember(request.Context, request.MemberId, h.summarizeRequestChannel)
			case request := <-*h.summarizeRequestChannel:
				go h.summarizer.Summarize(request, h.openPullRequestRequestChannel)
			case request := <-*h.openPullRequestRequestChannel:
				go h.pullRequestOpener.OpenPullRequest(request)
			}
		}
	}()
}
