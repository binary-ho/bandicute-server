package study_member

import (
	"bandicute-server/internal/storage/repository/connection"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, studyMemberID string) (*Model, error)
	GetAllByStudyId(ctx context.Context, studyID string) ([]*Model, error)
	GetAllMemberId(ctx context.Context) ([]string, error)
}

func NewStudyMemberRepository(base connection.Connection) Repository {
	return &StudyMemberRepository{Connection: base}
}
