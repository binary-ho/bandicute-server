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

func (r *SummaryRepository) CreatePostSummary(ctx context.Context, summary *summary.Model) error {
	req, err := r.NewRequest(ctx, PostMethod, "BaseEndpoint/post_summaries", summary)
	if err != nil {
		return err
	}

	return r.Do(req, summary)
}

func (r *SummaryRepository) GetPostSummary(ctx context.Context, blogPostID string) (*summary.Model, error) {
	req, err := r.NewRequest(ctx, GetMethod, "BaseEndpoint/post_summaries?blog_post_id=eq."+blogPostID, nil)
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

func (r *SummaryRepository) UpdatePostSummary(ctx context.Context, summary *summary.Model) error {
	req, err := r.NewRequest(ctx, PatchMethod, "BaseEndpoint/post_summaries?id=eq."+summary.ID, summary)
	if err != nil {
		return err
	}

	return r.Do(req, nil)
}
