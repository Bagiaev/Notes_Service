package main

import (
	"os"

	"notes_service/internal/service"
	"notes_service/pgk/jwt"
	"notes_service/pgk/logs"
	"notes_service/pgk/middleware"
	"notes_service/pgk/validator"

	"github.com/labstack/echo/v4"
)

func main() {
	// создаем логгер
	logger := logs.NewLogger(false)

	// подключаемся к базе
	db, err := PostgresConnection()
	if err != nil {
		logger.Fatal(err)
	}

	//jwt
	jwtSecret := jwt.NewJWT(os.Getenv("JWT_SECRET"))

	svc := service.NewService(db, logger, *jwtSecret)

	router := echo.New()
	//Валидатор
	router.Validator = validator.New()
	api := router.Group("api")

	//прописываем ручки
	api.GET("/note/:id", svc.GetNoteById)

	api.POST("/register", svc.Register)
	api.POST("/login", svc.Login)

	//защищеные роуты
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(jwtSecret))

	protected.GET("/profile", svc.ProfileHandler)

	//запуск сервера
	router.Logger.Fatal(router.Start(":8000"))

}
