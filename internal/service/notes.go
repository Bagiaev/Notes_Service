package service

import (
	"encoding/json"
	"io"
	"net/http"
	"notes_service/internal/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Получение цитат из сайта
func (s *Service) getInspirationalQuote() (string, error) {
	resp, err := http.Get("https://favqs.com/api/qotd")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Quote struct {
			Body string `json:"body"`
		} `json:"quote"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return "\n\nInspirational quote: " + result.Quote.Body, nil
}

// Функция для создания заметки
func (s *Service) CreateNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	var req models.NoteRequest

	if err := c.Bind(&req); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	if err := c.Validate(req); err != nil {
		s.logger.Error("Validation failed:", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	// Получаем цитату
	quote, err := s.getInspirationalQuote()
	if err != nil {
		s.logger.Error("Failed to get quote:", err)
		quote = "" // Продолжаем без цитаты при ошибке
	}

	note := &models.Note{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body + quote,
	}

	if err := s.notesRepo.RCreateNote(note); err != nil {
		s.logger.Error("Failed to create note:", err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusCreated, note)
}

// Функция для получения заметок юзера
func (s *Service) GetUserNotes(c echo.Context) error {
	userID := c.Get("userID").(int)

	notes, err := s.notesRepo.RGetUserNotes(userID)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: notes})
}

// Функция GetNoteById для получения заметки по id
//
//localhost:8000/api/note/:id
func (s *Service) GetNoteById(c echo.Context) error {
	userID := c.Get("userID").(int)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.notesRepo
	note, err := repo.RGetNoteById(id, userID)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: note})
}

// Функция для обновления заметок
func (s *Service) UpdateNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	var req models.NoteRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	if err := c.Validate(req); err != nil {
		s.logger.Error("Validation failed:", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	note := &models.Note{
		ID:     id,
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
	}

	if err := s.notesRepo.RUpdateNote(note); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, note)
}

// Функия для удаления заметок
func (s *Service) DeleteNote(c echo.Context) error {
	userID := c.Get("userID").(int)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	if err := s.notesRepo.RDeleteNote(id, userID); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.NoContent(http.StatusNoContent)
}
