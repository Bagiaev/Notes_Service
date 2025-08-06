package main

import (
	"os"

	"notes_service/internal/middleware"
	"notes_service/internal/service"
	"notes_service/pkg/jwt"
	"notes_service/pkg/logs"
	"notes_service/pkg/validator"

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
	api.POST("/register", svc.Register)
	api.POST("/login", svc.Login)

	//защищеные роуты
	protected := api.Group("user")
	protected.Use(middleware.JWTMiddleware(jwtSecret))

	protected.GET("/profile", svc.ProfileHandler) //Профиль
	//4 пункт
	protected.POST("/notes", svc.CreateNote)       //Создание заметки
	protected.GET("/notes", svc.GetUserNotes)      //Получение заметок пользователя
	protected.GET("/notes/:id", svc.GetNoteById)   //Получение заметки по id и по user_id
	protected.PUT("/notes/:id", svc.UpdateNote)    //Обновление заметки
	protected.DELETE("/notes/:id", svc.DeleteNote) //Удаление заметки по id  user_id

	//запуск сервера
	router.Logger.Fatal(router.Start(":8000"))

}
