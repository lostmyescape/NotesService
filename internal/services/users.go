package services

import (
	"NotesService/internal/models/users"
	"NotesService/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

var validate = validator.New()

// Register регистрация нового пользователя
func (s *Service) Register(c echo.Context) error {
	var user users.UserRequest

	if err := c.Bind(&user); err != nil {
		log.Printf("Ошибка передачи данных из контекста: %v", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	// валидация полей
	if err := validate.Var(user.Email, "required,email"); err != nil {
		return c.JSON(s.NewError("Поле email не соответствует минимальным требованиям"))
	}
	if err := validate.Var(user.Password, "required,min=6,max=30"); err != nil {
		return c.JSON(s.NewError("Минимальная длина пароля 8, а максимальная 30 символов"))
	}

	// хешируем пароль
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	repo := s.usersRepo

	// вносим данные в таблицу
	if err = repo.RegisterByEmail(user.Email, hashedPassword, time.Now()); err != nil {
		log.Printf("Ошибка при попытке выполнить sql-запрос: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	return c.JSON(http.StatusOK, "Пользователь был создан")
}

func (s *Service) Login(c echo.Context) error {
	var user users.UserRequest
	if err := c.Bind(&user); err != nil {
		log.Printf("Ошибка передачи данных из контекста: %v", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	// соответствие полей
	if err := validate.Struct(user); err != nil {
		return c.JSON(s.NewError("Ошибка авторизации"))
	}

	repo := s.usersRepo

	// достаем пароль
	storedPassword, err := repo.SelectPasswordByEmail(user.Email)
	if err != nil {
		log.Printf("Ошибка при попытке достать пароль по имейлу: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	// достаем айди
	userID, err := repo.SelectIdByEmail(user.Email)
	if err != nil {
		log.Printf("Ошибка при попытке достать айди по имейлу: %v", err)
		return c.JSON(s.NewError(InternalServerError))
	}

	// проверка паролей
	if err := utils.CheckPassword(storedPassword, user.Password); err != nil {
		log.Printf("Пароли не совпадают: %v", err)
		return c.JSON(s.NewError("Ошибка логина или пароля"))
	}

	// генерируем токен
	token, err := utils.GenerateToken(userID)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}
