package service

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Функция GetNoteById для получения заметки по id
//
//localhost:8000/api/note/:id
func (s *Service) GetNoteById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.notesRepo
	word, err := repo.RGetNoteById(id)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: word})
}
