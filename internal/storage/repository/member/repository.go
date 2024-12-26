package member

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	GetMember(ctx context.Context, id string) (*Model, error)
	GetMemberByBlog(ctx context.Context, blog string) (*Model, error)
}

func NewMemberRepository(conn connection.Connection) Repository {
	return &supabase.MemberRepository{Connection: conn}
}
