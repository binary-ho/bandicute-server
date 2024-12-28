package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/repository/post"
	"context"
	"fmt"
)

type PostRepository struct {
	Connection connection.Connection
}

func (r *PostRepository) GetById(ctx context.Context, id string) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, getMethod, post.TableName, "id=eq."+id, nil)
	if err != nil {
		return nil, err
	}

	var posts []*post.Model
	if err := r.Connection.Do(req, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("post post not found: %s", id)
	}

	return posts[0], nil
}

func (r *PostRepository) GetLatestByMemberId(ctx context.Context, memberID string) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, getMethod, post.TableName, "member_id=eq."+memberID+"&order=published_at.desc&limit=1", nil)
	if err != nil {
		return nil, err
	}

	var posts []*post.Model
	if err := r.Connection.Do(req, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("no post posts found for member: %s", memberID)
	}

	return posts[0], nil
}

func (r *PostRepository) GetByGUID(ctx context.Context, guid string) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, getMethod, post.TableName, "guid=eq."+guid, nil)
	if err != nil {
		return nil, err
	}

	var posts []*post.Model
	if err := r.Connection.Do(req, &posts); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("post post not found with guid: %s", guid)
	}

	return posts[0], nil
}

func (r *PostRepository) Create(ctx context.Context, model *post.Model) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, postMethod, post.TableName, "", *model)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, *model)
	return model, err
}

func (r *PostRepository) CreateAll(ctx context.Context, models []*post.Model) ([]*post.Model, error) {
	modelValues := make([]post.Model, len(models))
	for _, model := range models {
		modelValues = append(modelValues, *model)
	}

	req, err := r.Connection.NewRequest(ctx, postMethod, post.TableName, "", modelValues)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, modelValues)
	for i, modelValue := range modelValues {
		models[i] = &modelValue
	}
	return models, err
}

func (r *PostRepository) Update(ctx context.Context, model *post.Model) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, patchMethod, post.TableName, "id=eq."+model.ID, *model)
	if err != nil {
		return nil, err
	}

	err = r.Connection.Do(req, *model)
	return model, err
}
