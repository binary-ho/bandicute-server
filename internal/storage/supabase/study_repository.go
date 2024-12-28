package supabase

import (
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/repository/study"
	studyMember "bandicute-server/internal/storage/repository/study-member"
	"context"
	"fmt"
)

type StudyRepository struct {
	connection.Connection
}

func (r *StudyRepository) GetById(ctx context.Context, id string) (*study.Model, error) {
	req, err := r.NewRequest(ctx, getMethod, study.TableName, "id=eq."+id, nil)
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

func (r *StudyRepository) GetStudyMembers(ctx context.Context, studyID string) ([]*studyMember.Model, error) {
	req, err := r.NewRequest(ctx, getMethod, study.TableName, "study_id=eq."+studyID, nil)
	if err != nil {
		return nil, err
	}

	var members []*studyMember.Model
	if err := r.Do(req, &members); err != nil {
		return nil, err
	}

	return members, nil
}
