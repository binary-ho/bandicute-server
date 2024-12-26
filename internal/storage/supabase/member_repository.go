package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/repository/member"
	"context"
	"fmt"
)

type MemberRepository struct {
	Connection connection.Connection
}

func (r *MemberRepository) GetMember(ctx context.Context, id string) (*member.Model, error) {
	req, err := r.Connection.NewRequest(ctx, GetMethod, "/members?id=eq."+id, nil)
	if err != nil {
		return nil, err
	}

	var members []*member.Model
	if err := r.Connection.Do(req, &members); err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("member not found: %s", id)
	}

	return members[0], nil
}

func (r *MemberRepository) GetMemberByBlog(ctx context.Context, tistoryBlog string) (*member.Model, error) {
	req, err := r.Connection.NewRequest(ctx, GetMethod, "/members?tistory_blog=eq."+tistoryBlog, nil)
	if err != nil {
		return nil, err
	}

	var members []*member.Model
	if err := r.Connection.Do(req, &members); err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("member not found with tistory post: %s", tistoryBlog)
	}

	return members[0], nil
}
