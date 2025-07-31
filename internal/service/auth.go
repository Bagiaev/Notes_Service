package service

import (
	"database/sql"
	"net/http"
	"notes_service/internal/models"

	"github.com/labstack/echo/v4"
)

func (s *Service) Register(c echo.Context) error {
	var req models.AuthRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	//Валидация данных
	//Проверяет по структуре AuthRequest (в данном случае
	// валидность email и пароль не меньше 6 знаков)
	if err := c.Validate(req); err != nil {
		s.logger.Error("validation failed:", err)
		return c.JSON(s.NewError(InvalidParams))
	}
	//проверка сущетсвует ли пользователь с таким email
	repo := s.authRepo
	existingUser, err := repo.GetUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}
	if existingUser != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "user already exists"})
	}

	user := &models.User{
		Email: req.Email,
	}
	if err := user.HashPassword(req.Password); err != nil {
		s.logger.Error("failed to hash password:", err)
		return c.JSON(s.NewError(InternalServerError))
	}

	err = repo.CreateUser(user)
	if err != nil {
		s.logger.Error("failed to create user:", err)
		return c.JSON(s.NewError(InternalServerError))
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (s *Service) Login(c echo.Context) error {
	var req models.AuthRequest
	//Биндим данные
	err := c.Bind(&req)
	if err != nil {
		s.logger.Error("failed bind req:", err)
		return c.JSON(s.NewError(InvalidParams))
	}
	//Валидация данных
	err = c.Validate(req)
	if err != nil {
		s.logger.Error("validation failed", err)
		return c.JSON(s.NewError(InvalidParams))
	}

	//Поиск пользователя
	repo := s.authRepo
	user, err := repo.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("user not found:", req.Email)
			return c.JSON(s.NewError(InternalServerError))
		}
		s.logger.Error("database error", err)
		return c.JSON(s.NewError(InternalServerError))
	}

	//Проверка пароля
	if !user.CheckPassword(req.Password) {
		s.logger.Error("invalid password", user.Email)
		return c.JSON(s.NewError(InvalidParams))
	}
	//генерация токена
	token, err := s.jwt.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("failed to generate token:", err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"token":  token,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// j
func (s *Service) ProfileHandler(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user ID"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Profile data",
		"user": map[string]interface{}{
			"id": userID,
		},
	})
}
