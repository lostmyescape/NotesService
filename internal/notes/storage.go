package notes

import "time"

// полная структура
type Notes struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id,omitempty"`
	Title     string    `json:"title" validate:"required,max=20"`
	Body      string    `json:"body" validate:"required,min=10,max=80"`
	CreatedAt time.Time `json:"created_at"`
}
