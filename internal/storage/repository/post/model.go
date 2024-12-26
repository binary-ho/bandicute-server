package post

import "time"

// Model represents a post
type Model struct {
	ID          string    `json:"id"`
	MemberID    string    `json:"member_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Content     string    `json:"content"`
	GUID        string    `json:"guid"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
