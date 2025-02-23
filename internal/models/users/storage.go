package users

import "time"

// User для данных, которые уже хранятся в базе данных
type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRequest для данных, которые приходят при регистрации/логине
type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
