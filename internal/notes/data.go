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

// Функция RGetNoteById, чтобы взять из бд заметку по id
func (r *Repo) RGetNoteById(id, userID int) (*models.Note, error) {
	var note models.Note
	err := r.db.QueryRow(`SELECT id, title, body, created_at FROM notes WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// Функция RCreateNote, чтобы добавить в бд новую заметку
func (r *Repo) RCreateNote(note *models.Note) error {
	err := r.db.QueryRow(`INSERT INTO notes (user_id, title, body) VALUES ($1, $2, $3) RETURNING id, created_at`, note.UserID, note.Title, note.Body).
		Scan(&note.ID, &note.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Функция RGetUserNotes которая получает все заметки юзера по id
func (r *Repo) RGetUserNotes(userID int) ([]models.Note, error) {
	rows, err := r.db.Query(
		`SELECT id, title, body, created_at 
         FROM notes WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// Функция RUpdateNote для обновления заметки
func (r *Repo) RUpdateNote(note *models.Note) error {
	_, err := r.db.Exec(
		`UPDATE notes SET title = $1, body = $2 
         WHERE id = $3 AND user_id = $4`,
		note.Title, note.Body, note.ID, note.UserID,
	)
	return err
}

// Функция RDeleteNote для удаления заметки по id
func (r *Repo) RDeleteNote(id, userID int) error {
	_, err := r.db.Exec(
		`DELETE FROM notes WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return err
}
