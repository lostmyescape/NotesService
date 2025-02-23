package users

import (
	"database/sql"
	"time"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// RegisterByEmail создаст в таблице нового юзера
func (r *Repo) RegisterByEmail(email, hashedPassword string, createdAt time.Time) error {

	_, err := r.db.Exec(`
		INSERT INTO users (email, hashed_password, created_at) 
		VALUES ($1, $2, $3)
		`, email, hashedPassword, createdAt)

	if err != nil {
		return err
	}

	return nil
}

// SelectPasswordByEmail достанет пароль из бд по имейлу
func (r *Repo) SelectPasswordByEmail(email string) (string, error) {
	var password string
	err := r.db.QueryRow(`SELECT hashed_password FROM users WHERE email = $1`, email).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}
