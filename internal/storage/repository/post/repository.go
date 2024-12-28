package post

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Model, error)
	GetLatestByMemberId(ctx context.Context, memberID string) (*Model, error)
	GetByGUID(ctx context.Context, guid string) (*Model, error)
	Create(ctx context.Context, model *Model) (*Model, error)
	CreateAll(ctx context.Context, models []*Model) ([]*Model, error)
	Update(ctx context.Context, model *Model) (*Model, error)
}

func NewPostRepository(conn connection.Connection) Repository {
	return &supabase.PostRepository{Connection: conn}
}
