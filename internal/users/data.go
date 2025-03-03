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

// SelectPasswordByEmail достанет пароль
func (r *Repo) SelectPasswordByEmail(email string) (string, error) {
	var password string
	err := r.db.QueryRow(`SELECT hashed_password FROM users WHERE email = $1`, email).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

// SelectIdAndEmail достанет айди и имейл
func (r *Repo) SelectIdAndEmail(email string) (int, string, error) {
	var userID int
	var storedPassword string

	err := r.db.QueryRow(`
		SELECT id, hashed_password FROM users WHERE email = $1
	`, email).Scan(&userID, &storedPassword)

	if err != nil {
		return 0, "", err
	}

	return userID, storedPassword, nil
}
