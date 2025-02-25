package users

// UserRequest для данных, которые приходят при регистрации/логине
type UserRequest struct {
	Id       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
