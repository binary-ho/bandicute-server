package service

import (
	"bandicute-server/internal/service/channel"
)

type Dispatcher struct {
	parser                            *Parser
	summarizer                        *Summarizer
	pullRequestOpener                 *PullRequestOpener
	parsePostByMemberIdRequestChannel *channel.ParsePostByMemberIdRequest
	summarizeRequestChannel           *channel.SummarizeRequest
	openPullRequestRequestChannel     *channel.OpenPullRequestRequest
}

func NewDispatcher(
	parser *Parser,
	summarizer *Summarizer,
	pullRequestOpener *PullRequestOpener,
	parsePostByMemberIdRequestChannel *channel.ParsePostByMemberIdRequest,
	summarizeRequestChannel *channel.SummarizeRequest,
	openPullRequestRequestChannel *channel.OpenPullRequestRequest,
) *Dispatcher {
	return &Dispatcher{
		parser:                            parser,
		summarizer:                        summarizer,
		pullRequestOpener:                 pullRequestOpener,
		parsePostByMemberIdRequestChannel: parsePostByMemberIdRequestChannel,
		summarizeRequestChannel:           summarizeRequestChannel,
		openPullRequestRequestChannel:     openPullRequestRequestChannel,
	}
}

func (h *Dispatcher) Run() {
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
