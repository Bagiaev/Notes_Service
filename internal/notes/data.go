package notes

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

func (r *Repo) RGetNoteById(id int) (*models.Note, error) {
	var note models.Note
	err := r.db.QueryRow(`SELECT id, user_id, title, body, created_at FROM notes WHERE id = $1`, id).
		Scan(&note.ID, &note.UserID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &note, nil
}
