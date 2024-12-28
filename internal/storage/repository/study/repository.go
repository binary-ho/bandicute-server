package study

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Model, error)
}

func NewStudyRepository(base connection.Connection) Repository {
	return &supabase.StudyRepository{Connection: base}
}
