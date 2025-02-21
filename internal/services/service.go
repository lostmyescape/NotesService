package services

import (
	"NotesService/internal/handlers/notes"
	"NotesService/internal/handlers/users"
	"database/sql"
	"github.com/labstack/echo/v4"
)

const (
	InvalidParams       = "invalid params"
	InternalServerError = "internal error"
)

type Service struct {
	db        *sql.DB
	logger    echo.Logger
	usersRepo *users.Repo
	notesRepo *notes.Repo
}

func NewService(db *sql.DB, logger echo.Logger) *Service {
	svc := &Service{
		db:     db,
		logger: logger,
	}
	svc.initRepositories(db)

	return svc
}

func (s *Service) initRepositories(db *sql.DB) {
	s.usersRepo = users.NewRepo(db)
	s.notesRepo = notes.NewRepo(db)
}
