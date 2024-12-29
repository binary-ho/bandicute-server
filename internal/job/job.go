package job

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

func NewJobHandler(
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
				h.parser.ParseRecentPostByMember(request.Context, request.MemberId, h.summarizeRequestChannel)
			case request := <-*h.summarizeRequestChannel:
				h.summarizer.Summarize(request, h.openPullRequestRequestChannel)
			case request := <-*h.openPullRequestRequestChannel:
				h.pullRequestOpener.OpenPullRequest(request)
			}
		}
	}()
}
