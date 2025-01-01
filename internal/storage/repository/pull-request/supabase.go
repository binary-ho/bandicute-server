package pull_request

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type PullRequestRepository struct {
	Connection connection.Connection
}

func (r *PullRequestRepository) GetByPostIdAndStudyId(ctx context.Context, blogPostID string, studyID string) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.GetMethod, TableName, fmt.Sprintf("blog_post_id=eq.%s&study_id=eq.%s", blogPostID, studyID), nil)
	if err != nil {
		return nil, err
	}

	var prs []*Model
	if err := r.Connection.Do(req, &prs); err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, fmt.Errorf("pull request not found for post post %s and study %s", blogPostID, studyID)
	}

	return prs[0], nil
}

func (r *PullRequestRepository) Create(ctx context.Context, model *Model) (*Model, error) {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}

	req, err := r.Connection.NewRequest(ctx, supabase.PostMethod, TableName, "", model)
	if err != nil {
		return nil, err
	}

	var result []Model
	err = r.Connection.Do(req, &result)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		*model = result[0]
	}
	return model, nil
}

func (r *PullRequestRepository) Update(ctx context.Context, pr *Model) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.PatchMethod, TableName, "id=eq."+pr.ID, pr)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, pr)
	return pr, err
}
