package services

import (
	"NotesService/internal/models/users"
	"NotesService/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var validate = validator.New()

// Register регистрация нового пользователя
func (s *Service) Register(c echo.Context) error {
	var user users.UserRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(s.NewError(InvalidParams))
	}

	// соответствие полей
	if err := validate.Struct(user); err != nil {
		return c.JSON(s.NewError("Ошибка регистрации"))
	}

	// хешируем пароль
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.JSON(s.NewError(InternalServerError))
	}

	repo := s.usersRepo

	// регистрация
	if err = repo.RegisterByEmail(user.Email, hashedPassword, time.Now()); err != nil {
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, "Пользователь был создан")
}

func (s *Service) Login(c echo.Context) error {
	var user users.UserRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(s.NewError(InvalidParams))
	}

	// соответствие полей
	if err := validate.Struct(user); err != nil {
		return c.JSON(s.NewError("Ошибка авторизации"))
	}

	repo := s.usersRepo

	// достаем пароль из бд
	storedPassword, err := repo.SelectPasswordByEmail(user.Email)
	if err != nil {
		return c.JSON(s.NewError(InternalServerError))
	}

	// проверка паролей
	if err := utils.CheckPassword(storedPassword, user.Password); err != nil {
		return c.JSON(s.NewError("Ошибка логина или пароля"))
	}

	// генерация токена
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}
