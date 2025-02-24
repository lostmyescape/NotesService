package notes

import "time"

// полная структура
type Notes struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Title     string    `json:"title" validate:"required,max=20"`
	Body      string    `json:"body" validate:"required,min=10,max=80"`
	CreatedAt time.Time `json:"created_at"`
}

// структура для получения заметок

type GetNotes struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
