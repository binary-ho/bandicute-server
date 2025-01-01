package member

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
	"fmt"
)

type MemberRepository struct {
	Connection connection.Connection
}

func (r *MemberRepository) GetById(ctx context.Context, id string) (*Model, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.GetMethod, TableName, "id=eq."+id, nil)
	if err != nil {
		return nil, err
	}

	var members []*Model
	if err := r.Connection.Do(req, &members); err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("member not found: %s", id)
	}

	return members[0], nil
}

func (r *MemberRepository) GetBlogUrlById(ctx context.Context, id string) (string, error) {
	req, err := r.Connection.NewRequest(ctx, supabase.GetMethod, TableName, "select=blog_url&id=eq."+id, nil)
	if err != nil {
		return "", err
	}

	var blogUrl string
	if err := r.Connection.Do(req, &blogUrl); err != nil {
		return "", err
	}

	return blogUrl, err
}
