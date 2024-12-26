package summary

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	CreatePostSummary(ctx context.Context, summary *Model) error
	GetPostSummary(ctx context.Context, blogPostID string) (*Model, error)
	UpdatePostSummary(ctx context.Context, summary *Model) error
}

func NewPostWriterRepository(base connection.Connection) Repository {
	return &supabase.SummaryRepository{Connection: base}
}
