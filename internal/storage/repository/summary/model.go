package summary

import (
	"bandicute-server/internal/storage/repository/connection"
	"time"
)

// Model represents a post
type Model struct {
	ID           string    `json:"id"`
	BlogPostID   string    `json:"blog_post_id"`
	Summary      string    `json:"summary"`
	IsSummarized bool      `json:"is_summarized"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const TableName = connection.Table("post_summaries")
