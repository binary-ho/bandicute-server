package summary

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetByPostId(ctx context.Context, blogPostID string) (*Model, error)
	Create(ctx context.Context, summary *Model) (*Model, error)
	Update(ctx context.Context, summary *Model) (*Model, error)
}

func NewPostWriterRepository(base connection.Connection) Repository {
	return &supabase.SummaryRepository{Connection: base}
}
