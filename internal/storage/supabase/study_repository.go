package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/repository/study"
	study_member "bandicute-server/internal/storage/repository/study-member"
	"context"
	"fmt"
)

type StudyRepository struct {
	connection.Connection
}

func (r *StudyRepository) GetStudy(ctx context.Context, id string) (*study.Model, error) {
	req, err := r.NewRequest(ctx, GetMethod, "BaseEndpoint/studies?id=eq."+id, nil)
	if err != nil {
		return nil, err
	}

	var studies []*study.Model
	if err := r.Do(req, &studies); err != nil {
		return nil, err
	}

	if len(studies) == 0 {
		return nil, fmt.Errorf("study not found: %s", id)
	}

	return studies[0], nil
}

func (r *StudyRepository) GetStudyMembers(ctx context.Context, studyID string) ([]*study_member.Model, error) {
	req, err := r.NewRequest(ctx, GetMethod, "BaseEndpoint/study_members?study_id=eq."+studyID, nil)
	if err != nil {
		return nil, err
	}

	var members []*study_member.Model
	if err := r.Do(req, &members); err != nil {
		return nil, err
	}

	return members, nil
}
