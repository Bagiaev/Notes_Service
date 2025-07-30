package main

import (
	"github.com/labstack/echo/v4"

	"notes_service/internal/service"
	"notes_service/pgk/logs"
)

func main() {
	// создаем логгер
	logger := logs.NewLogger(false)

	// подключаемся к базе
	db, err := PostgresConnection()
	if err != nil {
		logger.Fatal(err)
	}

	svc := service.NewService(db, logger)

	router := echo.New()
	api := router.Group("api")

	//прописываем ручки
	api.GET("/note/:id", svc.GetNoteById)

	//запуск сервера
	router.Logger.Fatal(router.Start(":8000"))

}
