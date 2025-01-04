package study_member

import (
	"bandicute-server/internal/storage/repository/connection"
	"time"
)

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

type MemberId struct {
	MemberId string `json:"member_id"`
}

const TableName = connection.Table("study_members")
const MemberIdView = connection.Table("study_members_member_id")
