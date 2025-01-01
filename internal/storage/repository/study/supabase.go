package study

import (
	"bandicute-server/internal/storage/repository/connection"
	studyMember "bandicute-server/internal/storage/repository/study-member"
	"bandicute-server/internal/storage/supabase"
	"context"
	"fmt"
)

type StudyRepository struct {
	connection.Connection
}

func (r *StudyRepository) GetById(ctx context.Context, id string) (*Model, error) {
	req, err := r.NewRequest(ctx, supabase.GetMethod, TableName, "id=eq."+id, nil)
	if err != nil {
		return nil, err
	}

	var studies []*Model
	if err := r.Do(req, &studies); err != nil {
		return nil, err
	}

	if len(studies) == 0 {
		return nil, fmt.Errorf("study not found: %s", id)
	}

	return studies[0], nil
}

func (r *StudyRepository) GetAllByMemberId(ctx context.Context, memberID string) ([]*Model, error) {
	query := "select=*" +
		",study_members!inner(member_id)" +
		"&study_members.member_id=eq." + memberID
	req, err := r.NewRequest(ctx, supabase.GetMethod, TableName, query, nil)
	if err != nil {
		return nil, err
	}

	var studies []*Model
	if err := r.Do(req, &studies); err != nil {
		return nil, err
	}

	return studies, nil
}

func (r *StudyRepository) GetStudyMembers(ctx context.Context, studyID string) ([]*studyMember.Model, error) {
	req, err := r.NewRequest(ctx, supabase.GetMethod, TableName, "study_id=eq."+studyID, nil)
	if err != nil {
		return nil, err
	}

	var members []*studyMember.Model
	if err := r.Do(req, &members); err != nil {
		return nil, err
	}

	return members, nil
}
