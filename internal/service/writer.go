package service

import (
	"bandicute-server/internal/service/channel"
	"bandicute-server/internal/service/request"
	studyMember "bandicute-server/internal/storage/repository/study-member"
	"bandicute-server/pkg/logger"
	"context"
)

type Writer struct {
	studyMemberRepository studyMember.Repository
	parseRequestChannel   *channel.ParsePostByMemberIdRequest
}

func NewWriter(studyMemberRepository studyMember.Repository, parseRequestChannel *channel.ParsePostByMemberIdRequest) *Writer {
	return &Writer{
		studyMemberRepository: studyMemberRepository,
		parseRequestChannel:   parseRequestChannel,
	}
}

// TODO: PageNation으로 일종의 Chunk-Oriented 흉내낼 수 있으면 좋을듯.
func (w *Writer) WriteAllMembersPost(ctx context.Context) {
	studyMemberIds, err := w.studyMemberRepository.GetAllMemberId(ctx)
	if err != nil {
		logger.Error("Failed to get study members", logger.Fields{
			"error": err.Error(),
		})
		return
	}

	for _, id := range studyMemberIds {
		*w.parseRequestChannel <- request.ParsePostByMemberId{
			Context:  ctx,
			MemberId: id,
		}
	}
}

func (w *Writer) WriteByStudy(ctx context.Context, studyId string) {
	studyMembers, err := w.studyMemberRepository.GetAllByStudyId(ctx, studyId)
	if err != nil {
		logger.Error("Failed to get study members", logger.Fields{
			"studyId": studyId,
			"error":   err.Error(),
		})
		return
	}

	for _, member := range studyMembers {
		*w.parseRequestChannel <- request.ParsePostByMemberId{
			Context:  ctx,
			MemberId: member.MemberID,
		}
	}
}

func (w *Writer) WriteByMember(ctx context.Context, memberId string) {
	*w.parseRequestChannel <- request.ParsePostByMemberId{
		Context:  ctx,
		MemberId: memberId,
	}
}
