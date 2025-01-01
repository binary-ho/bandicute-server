package summary

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type SummaryRepository struct {
	connection.Connection
}

func (r *SummaryRepository) GetByPostId(ctx context.Context, blogPostID string) (*Model, error) {
	req, err := r.NewRequest(ctx, supabase.GetMethod, TableName, "blog_post_id=eq."+blogPostID, nil)
	if err != nil {
		return nil, err
	}

	var summaries []*Model
	if err := r.Do(req, &summaries); err != nil {
		return nil, err
	}

	if len(summaries) == 0 {
		return nil, fmt.Errorf("post summary not found for post post: %s", blogPostID)
	}

	return summaries[0], nil
}

func (r *SummaryRepository) Create(ctx context.Context, model *Model) (*Model, error) {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}

	req, err := r.NewRequest(ctx, supabase.PostMethod, TableName, "", model)
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

func (r *SummaryRepository) CreateAll(ctx context.Context, models []*Model) ([]*Model, error) {
	modelValues := make([]Model, len(models))
	for _, model := range models {
		modelValues = append(modelValues, *model)
	}

	req, err := r.Connection.NewRequest(ctx, supabase.PostMethod, TableName, "", modelValues)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, modelValues)
	for i, modelValue := range modelValues {
		models[i] = &modelValue
	}
	return models, err
}

func (r *SummaryRepository) Update(ctx context.Context, model *Model) (*Model, error) {
	req, err := r.NewRequest(ctx, supabase.PatchMethod, TableName, "id=eq."+model.ID, model)
	if err != nil {
		return nil, err
	}

	err = r.Do(req, model)
	return model, err
}
