package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	studyMember "bandicute-server/internal/storage/repository/study-member"
	"context"
)

type StudyMemberRepository struct {
	connection.Connection
}

func (r *StudyMemberRepository) GetById(ctx context.Context, studyMemberID string) (*studyMember.Model, error) {
	req, err := r.NewRequest(ctx, getMethod, studyMember.TableName, "id=eq."+studyMemberID, nil)
	if err != nil {
		return nil, err
	}

	var member *studyMember.Model
	if err := r.Do(req, &member); err != nil {
		return nil, err
	}

	return member, nil
}

func (r *StudyMemberRepository) GetAllByStudyId(ctx context.Context, studyID string) ([]*studyMember.Model, error) {
	req, err := r.NewRequest(ctx, getMethod, studyMember.TableName, "study_id=eq."+studyID, nil)
	if err != nil {
		return nil, err
	}

	var members []*studyMember.Model
	if err := r.Do(req, &members); err != nil {
		return nil, err
	}

	return members, nil
}
