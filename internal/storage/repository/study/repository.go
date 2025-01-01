package study

import (
	"bandicute-server/internal/storage/repository/connection"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Model, error)
	GetAllByMemberId(ctx context.Context, memberID string) ([]*Model, error)
}

func NewStudyRepository(base connection.Connection) Repository {
	return &StudyRepository{Connection: base}
}
