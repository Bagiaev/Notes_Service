package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email" validate:"required,email"`
	HashedPassword string    `json:"-" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (u *User) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
