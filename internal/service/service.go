package service

import (
	"database/sql"
	"notes_service/internal/auth"
	"notes_service/internal/notes"
	"notes_service/pkg/jwt"
	"notes_service/pkg/validator"

	"github.com/labstack/echo"
)

const (
	InvalidParams       = "invalid params"
	InternalServerError = "internal error"
)

type Service struct {
	db        *sql.DB
	logger    echo.Logger
	validator *validator.CustomValidator

	notesRepo *notes.Repo
	authRepo  *auth.Repo

	jwt *jwt.JWT
}

func NewService(db *sql.DB, logger echo.Logger, jwtSecret jwt.JWT) *Service {
	svc := &Service{
		db:        db,
		logger:    logger,
		validator: validator.New(),
		jwt:       &jwtSecret,
	}

	svc.initRepositories(db)

	return svc
}

func (s *Service) initRepositories(db *sql.DB) {
	s.notesRepo = notes.NewRepo(db)
	s.authRepo = auth.NewRepo(db)
}

// a
type Response struct {
	Object       any    `json:"object,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (r *Response) Error() string {
	return r.ErrorMessage
}

func (s *Service) NewError(err string) (int, *Response) {
	return 400, &Response{ErrorMessage: err}
}
