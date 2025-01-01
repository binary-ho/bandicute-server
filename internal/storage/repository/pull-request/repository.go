package pull_request

import (
	"bandicute-server/internal/storage/repository/connection"
	"context"
)

type Repository interface {
	GetByPostIdAndStudyId(ctx context.Context, blogPostID string, studyID string) (*Model, error)
	Create(ctx context.Context, pr *Model) (*Model, error)
	Update(ctx context.Context, pr *Model) (*Model, error)
}

func NewPullRequestRepository(conn connection.Connection) Repository {
	return &PullRequestRepository{Connection: conn}
}
