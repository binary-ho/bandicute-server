package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	pull_request "bandicute-server/internal/storage/repository/pull-request"
	"context"
	"fmt"
)

type PullRequestRepository struct {
	Connection connection.Connection
}

func (r *PullRequestRepository) CreatePullRequest(ctx context.Context, pr *pull_request.Model) error {
	req, err := r.Connection.NewRequest(ctx, PostMethod, "BaseEndpoint/pull_requests", pr)
	if err != nil {
		return err
	}

	return r.Connection.Do(req, pr)
}

func (r *PullRequestRepository) GetPullRequest(ctx context.Context, blogPostID string, studyID string) (*pull_request.Model, error) {
	req, err := r.Connection.NewRequest(ctx, GetMethod, fmt.Sprintf("BaseEndpoint/pull_requests?blog_post_id=eq.%s&study_id=eq.%s", blogPostID, studyID), nil)
	if err != nil {
		return nil, err
	}

	var prs []*pull_request.Model
	if err := r.Connection.Do(req, &prs); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, fmt.Errorf("pull request not found for post post %s and study %s", blogPostID, studyID)
	}

	return prs[0], nil
}

func (r *PullRequestRepository) UpdatePullRequest(ctx context.Context, pr *pull_request.Model) error {
	req, err := r.Connection.NewRequest(ctx, PatchMethod, "BaseEndpoint/pull_requests?id=eq."+pr.ID, pr)
	if err != nil {
		return err
	}

	return r.Connection.Do(req, nil)
}
