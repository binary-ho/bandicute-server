package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/repository/summary"
	"context"
	"fmt"
)

type SummaryRepository struct {
	connection.Connection
}

func (r *SummaryRepository) GetByPostId(ctx context.Context, blogPostID string) (*summary.Model, error) {
	req, err := r.NewRequest(ctx, getMethod, summary.TableName, "blog_post_id=eq."+blogPostID, nil)
	if err != nil {
		return nil, err
	}

	var summaries []*summary.Model
	if err := r.Do(req, &summaries); err != nil {
		return nil, err
	}

	if len(summaries) == 0 {
		return nil, fmt.Errorf("post summary not found for post post: %s", blogPostID)
	}

	return summaries[0], nil
}

func (r *SummaryRepository) Create(ctx context.Context, model *summary.Model) (*summary.Model, error) {
	req, err := r.NewRequest(ctx, postMethod, summary.TableName, "", model)
	if err != nil {
		return nil, err
	}

	err = r.Do(req, model)
	return model, err
}

func (r *SummaryRepository) Update(ctx context.Context, model *summary.Model) (*summary.Model, error) {
	req, err := r.NewRequest(ctx, patchMethod, summary.TableName, "id=eq."+model.ID, model)
	if err != nil {
		return nil, err
	}

	err = r.Do(req, model)
	return model, err
}
