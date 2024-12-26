package study_member

import "time"

// Model represents a post
type Model struct {
	ID        string    `json:"id"`
	StudyId   string    `json:"study_id"`
	MemberID  string    `json:"member_id"`
	Folder    string    `json:"folder_path"`
	IsLeader  bool      `json:"is_leader"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
