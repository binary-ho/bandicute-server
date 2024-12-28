package study

import (
	"bandicute-server/internal/storage/repository/connection"
	"time"
)

// Model represents a post
type Model struct {
	ID          string    `json:"id"`
	MemberID    string    `json:"member_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Repository  string    `json:"github_repo"`
	Branch      string    `json:"branch"`
	Directory   string    `json:"directory"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const TableName = connection.Table("studies")
