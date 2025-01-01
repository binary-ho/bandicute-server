package member

import (
	"bandicute-server/internal/storage/repository/connection"
	"time"
)

// Model represents a study member
type Model struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Blog      string    `json:"blog_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const TableName = connection.Table("members")
