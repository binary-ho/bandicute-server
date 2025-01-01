package post

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type PostRepository struct {
	Connection connection.Connection
}

func (r *PostRepository) GetById(ctx context.Context, id string) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.GetMethod, TableName, "id=eq."+id, nil)
	if err != nil {
		return nil, err
	}

	var posts []*Model
	if err := r.Connection.Do(req, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("post post not found: %s", id)
	}

	return posts[0], nil
}

func (r *PostRepository) GetLatestByMemberId(ctx context.Context, memberID string) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.GetMethod, TableName, "member_id=eq."+memberID+"&order=published_at.desc&limit=1", nil)
	if err != nil {
		return nil, err
	}

	var posts []*Model
	if err := r.Connection.Do(req, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no post posts found for member: %s", memberID)
	}

	return posts[0], nil
}

func (r *PostRepository) GetByGUID(ctx context.Context, guid string) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.GetMethod, TableName, "guid=eq."+guid, nil)
	if err != nil {
		return nil, err
	}

	var posts []*Model
	if err := r.Connection.Do(req, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("post post not found with guid: %s", guid)
	}

	return posts[0], nil
}

func (r *PostRepository) Create(ctx context.Context, model *Model) (*Model, error) {
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

func (r *PostRepository) CreateAll(ctx context.Context, models []*Model) ([]*Model, error) {
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

func (r *PostRepository) Update(ctx context.Context, model *Model) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.PatchMethod, TableName, "id=eq."+model.ID, model)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, model)
	return model, err
}
