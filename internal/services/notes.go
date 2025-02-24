package services

import (
	notes2 "NotesService/internal/models/notes"
	"NotesService/internal/services/quotes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

// CreateNoteHandler создает новую записку в бд
func (s *Service) CreateNoteHandler(c echo.Context) error {
	// получаем айди юзера из контекста
	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var note notes2.Notes

	if err := c.Bind(&note); err != nil {
		log.Printf("Ошибка передачи данных из контекста: %v", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	// валидация полей
	if err := validate.Var(note.Title, "required,max=20"); err != nil {
		return c.JSON(s.NewError("Максимальная длина тайтла 20 символов"))
	}
	if err := validate.Var(note.Body, "required,min=10,max=80"); err != nil {
		return c.JSON(s.NewError("Минимальная длина бади 10, а максимальная 80 символов"))
	}

	// получаем цитату из внешнего апи
	quote, err := quotes.GetQuote()
	if err != nil {
		log.Printf("Ошибка при работе с внешнем апи: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	// в конец заметки дополняем цитатой, взятой из внешнего апи
	note.Body = fmt.Sprintf("%s, цитата: %s", note.Body, quote)

	repo := s.notesRepo

	// создаем записку
	if err := repo.CreateNote(userID, note.Title, note.Body, time.Now()); err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	return c.JSON(http.StatusOK, "Метка была добавлена")
}

func (s *Service) GetAllNotesHandler(c echo.Context) error {
	// получаем айди юзера из контекста
	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	repo := s.notesRepo

	notes, err := repo.GetAllNotes(userID)
	if err != nil {
		log.Printf("Ошибка выполнения sql-запроса: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	return c.JSON(http.StatusOK, notes)

}
