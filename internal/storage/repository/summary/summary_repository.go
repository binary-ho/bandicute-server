package summary

import (
	"bandicute-server/internal/storage/repository/connection"
	"context"
)

type Repository interface {
	GetByPostId(ctx context.Context, blogPostID string) (*Model, error)
	Create(ctx context.Context, summary *Model) (*Model, error)
	CreateAll(ctx context.Context, summaries []*Model) ([]*Model, error)
	Update(ctx context.Context, summary *Model) (*Model, error)
}

func NewPostWriterRepository(base connection.Connection) Repository {
	return &SummaryRepository{Connection: base}
}
