package auth

import (
	"database/sql"
	"notes_service/internal/models"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(user *models.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (email, hashed_password, created_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING id, created_at`,
		user.Email,
		user.HashedPassword,
	)
	if err != nil {
		return err
	}
	return nil

}

func (r *Repo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`SELECT id, email, hashed_password, created_at FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
