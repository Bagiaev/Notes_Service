package models

import "time"

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title" validate:"required"`
	Body      string    `json:"body" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}
