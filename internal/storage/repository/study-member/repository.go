package study_member

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetAllByStudyId(ctx context.Context, studyID string) ([]*Model, error)
	GetById(ctx context.Context, studyMemberID string) (*Model, error)
}

func NewStudyMemberRepository(base connection.Connection) Repository {
	return &supabase.StudyMemberRepository{Connection: base}
}
