package member

import (
	"bandicute-server/internal/storage/repository/connection"
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*Model, error)
	GetBlogUrlById(ctx context.Context, id string) (string, error)
}

func NewMemberRepository(conn connection.Connection) Repository {
	return &MemberRepository{Connection: conn}
}
