package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	studymember "bandicute-server/internal/storage/repository/study-member"
	"context"
)

type StudyMemberRepository struct {
	connection.Connection
}

func (r *StudyMemberRepository) GetStudyMembers(ctx context.Context, studyID string) ([]*studymember.Model, error) {
	req, err := r.NewRequest(ctx, GetMethod, "/study_members?study_id=eq."+studyID, nil)
	if err != nil {
		return nil, err
	}

	var members []*studymember.Model
	if err := r.Do(req, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *StudyMemberRepository) GetStudyMember(ctx context.Context, studyMemberID string) (*studymember.Model, error) {
	req, err := r.NewRequest(ctx, GetMethod, "/study_members?id=eq."+studyMemberID, nil)
	if err != nil {
		return nil, err
	}

	var member *studymember.Model
	if err := r.Do(req, &member); err != nil {
		return nil, err
	}

	return member, nil
}
