package post

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	CreatePost(ctx context.Context, post *Model) error
	GetLatestPost(ctx context.Context, memberID string) (*Model, error)
	GetPost(ctx context.Context, id string) (*Model, error)
	GetPostByGUID(ctx context.Context, guid string) (*Model, error)
	UpdatePost(ctx context.Context, post *Model) error
}

func NewPostRepository(conn connection.Connection) Repository {
	return &supabase.PostRepository{Connection: conn}
}
