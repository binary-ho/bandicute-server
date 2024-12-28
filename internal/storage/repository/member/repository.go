package member

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Model, error)
	GetBlogUrlById(ctx context.Context, id string) (string, error)
}

func NewMemberRepository(conn connection.Connection) Repository {
	return &supabase.MemberRepository{Connection: conn}
}
