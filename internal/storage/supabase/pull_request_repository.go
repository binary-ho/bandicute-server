package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	pullRequest "bandicute-server/internal/storage/repository/pull-request"
	"context"
	"fmt"
)

type PullRequestRepository struct {
	Connection connection.Connection
}

func (r *PullRequestRepository) GetByPostIdAndStudyId(ctx context.Context, blogPostID string, studyID string) (*pullRequest.Model, error) {
	req, err := r.Connection.NewRequest(ctx, getMethod, pullRequest.TableName, fmt.Sprintf("blog_post_id=eq.%s&study_id=eq.%s", blogPostID, studyID), nil)
	if err != nil {
		return nil, err
	}

	var prs []*pullRequest.Model
	if err := r.Connection.Do(req, &prs); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, fmt.Errorf("pull request not found for post post %s and study %s", blogPostID, studyID)
	}

	return prs[0], nil
}

func (r *PullRequestRepository) Create(ctx context.Context, pr *pullRequest.Model) (*pullRequest.Model, error) {
	req, err := r.Connection.NewRequest(ctx, postMethod, pullRequest.TableName, "", pr)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, pr)
	return pr, err
}

func (r *PullRequestRepository) Update(ctx context.Context, pr *pullRequest.Model) (*pullRequest.Model, error) {
	req, err := r.Connection.NewRequest(ctx, patchMethod, pullRequest.TableName, "id=eq."+pr.ID, pr)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, pr)
	return pr, err
}
