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
func (r *Repo) GetAllNotes(userID int) (*[]notes2.GetNotes, error) {
	var notes []notes2.GetNotes

	query := `
		SELECT id, title, body, created_at 
		FROM notes
		WHERE user_id = $1
		`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var note notes2.GetNotes
		err := rows.Scan(&note.Id, &note.Title, &note.Body, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return &notes, nil
}

// GetNoteById получаем записку у авторизированного юзера по айди заметки
func (r *Repo) GetNoteById(userID, id int) (*notes2.GetNotes, error) {
	var note notes2.GetNotes
	err := r.db.QueryRow(`
								SELECT id, title, body, created_at 
								FROM notes 
								WHERE user_id = $1 AND id = $2`, userID, id).
		Scan(&note.Id, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *Repo) UpdateNoteById(title, body string, userID, id int) error {

	_, err := r.db.Exec(`
								UPDATE notes						
								SET title = $1, body = $2							
								WHERE user_id = $3 AND id = $4
						`, title, body, userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteNoteById(userID, id int) error {

	_, err := r.db.Exec(`DELETE FROM notes WHERE id = $1 AND user_id = $2`, id, userID)

	if err != nil {
		return err
	}

	return nil
}
