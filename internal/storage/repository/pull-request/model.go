package pull_request

import "time"

// Model represents a GitHub pull request
type Model struct {
	ID         string    `json:"id"`
	BlogPostID string    `json:"blog_post_id"`
	StudyID    string    `json:"study_id"`
	PrUrl      string    `json:"pr_url"`
	IsOpened   bool      `json:"is_opened"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
