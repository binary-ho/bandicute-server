package pull_request

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type Repository interface {
	CreatePullRequest(ctx context.Context, pr *Model) error
	GetPullRequest(ctx context.Context, blogPostID string, studyID string) (*Model, error)
	UpdatePullRequest(ctx context.Context, pr *Model) error
}

func NewPullRequestRepository(conn connection.Connection) Repository {
	return &supabase.PullRequestRepository{Connection: conn}
}
