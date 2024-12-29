package study_member

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, studyMemberID string) (*Model, error)
	GetAllByStudyId(ctx context.Context, studyID string) ([]*Model, error)
}

func NewStudyMemberRepository(base connection.Connection) Repository {
	return &supabase.StudyMemberRepository{Connection: base}
}
