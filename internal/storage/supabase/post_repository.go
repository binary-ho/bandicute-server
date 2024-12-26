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

func (r *PostRepository) CreatePost(ctx context.Context, post *post.Model) error {
	req, err := r.Connection.NewRequest(ctx, PostMethod, "/blog_posts", post)
	if err != nil {
		return err
	}

	return r.Connection.Do(req, post)
}

func (r *PostRepository) GetLatestPost(ctx context.Context, memberID string) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, GetMethod, fmt.Sprintf("blog_posts?member_id=eq.%s&order=published_at.desc&limit=1", memberID), nil)
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

func (r *PostRepository) GetPost(ctx context.Context, id string) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, GetMethod, "/blog_posts?id=eq."+id, nil)
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

func (r *PostRepository) GetPostByGUID(ctx context.Context, guid string) (*post.Model, error) {
	req, err := r.Connection.NewRequest(ctx, GetMethod, "/blog_posts?guid=eq."+guid, nil)
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

func (r *PostRepository) UpdatePost(ctx context.Context, post *post.Model) error {
	req, err := r.Connection.NewRequest(ctx, PatchMethod, "/blog_posts?id=eq."+post.ID, post)
	if err != nil {
		return err
	}

	return r.Connection.Do(req, nil)
}
