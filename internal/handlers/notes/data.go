package notes

import (
	notes2 "NotesService/internal/models/notes"
	"database/sql"
	"time"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// CreateNote добавит новую заметку в бд
func (r *Repo) CreateNote(userID int, title, body string, createdAt time.Time) error {
	_, err := r.db.Exec(`
							INSERT INTO notes (user_id, title, body, created_at)
    						VALUES ($1, $2, $3, $4)`, userID, title, body, createdAt)
	if err != nil {
		return err
	}

	return nil
}

// GetAllNotes получит все записки авторизированного юзера
func (r *Repo) GetAllNotes(userID int) ([]notes2.GetNotes, error) {
	var notes []notes2.GetNotes

	query := `
		SELECT title, body, created_at 
		FROM notes
		WHERE user_id = $1
		`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var note notes2.GetNotes
		err := rows.Scan(&note.Title, &note.Body, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return notes, nil

}
