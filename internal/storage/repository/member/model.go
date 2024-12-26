package member

import "time"

// Model represents a study member
type Model struct {
	ID        string    `json:"id"`
	Blog      string    `json:"blog"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
