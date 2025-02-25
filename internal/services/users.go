package services

import (
	"NotesService/internal/users"
	"NotesService/pkg/jwt"
	"NotesService/pkg/password"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Поле email не соответствует минимальным требованиям"})
	}
	if err := validate.Var(user.Password, "required,min=6"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Минимальная длина пароля 6 символов"})
	}

	// хешируем пароль
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	repo := s.usersRepo

	// вносим данные в таблицу
	if err = repo.RegisterByEmail(user.Email, hashedPassword, time.Now()); err != nil {
		// отлов ошибки с идентичным имейлом
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Пользователь с таким email уже существует"})
		}
		log.Printf("Ошибка при попытке выполнить sql-запрос: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "Пользователь был успешно создан"})
}

func (s *Service) Login(c echo.Context) error {
	var user users.UserRequest
	if err := c.Bind(&user); err != nil {
		log.Printf("Ошибка передачи данных из контекста: %v", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	// соответствие полей
	if err := validate.Struct(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Ошибка авторизации"})
	}

	repo := s.usersRepo

	// достаем айди и пароль
	userID, storedPassword, err := repo.SelectIdAndEmail(user.Email)
	if err != nil {
		log.Printf("Ошибка при попытке достать айди по имейлу: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неправильный логин или пароль"})
	}

	// проверка паролей
	if err := password.CheckPassword(storedPassword, user.Password); err != nil {
		log.Printf("Пароли не совпадают: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неправильный логин или пароль"})
	}

	// генерируем токен
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}
