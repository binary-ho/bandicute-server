package study_member

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/supabase"
	"context"
)

type StudyMemberRepository struct {
	connection.Connection
}

func (r *StudyMemberRepository) GetById(ctx context.Context, studyMemberID string) (*Model, error) {
	req, err := r.NewRequest(ctx, supabase.GetMethod, TableName, "id=eq."+studyMemberID, nil)
	if err != nil {
		return nil, err
	}

	var member *Model
	if err := r.Do(req, &member); err != nil {
		return nil, err
	}

	return member, nil
}

func (r *StudyMemberRepository) GetAllByStudyId(ctx context.Context, studyID string) ([]*Model, error) {
	req, err := r.NewRequest(ctx, supabase.GetMethod, TableName, "study_id=eq."+studyID, nil)
	if err != nil {
		return nil, err
	}

	members := make([]*Model, 0)
	if err := r.Do(req, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *StudyMemberRepository) GetAllMemberId(ctx context.Context) ([]string, error) {
	req, err := r.NewRequest(ctx, supabase.GetMethod, MemberIdView, "select=member_id", nil)
	if err != nil {
		return nil, err
	}

	memberIds := make([]*MemberId, 0)
	if err := r.Do(req, &memberIds); err != nil {
		return nil, err
	}
	return convertToStrings(memberIds), nil
}

func convertToStrings(memberIds []*MemberId) []string {
	stringIds := make([]string, len(memberIds))
	for _, id := range memberIds {
		stringIds = append(stringIds, id.MemberId)
	}
	return stringIds
}
