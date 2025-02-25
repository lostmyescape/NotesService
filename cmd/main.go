package main

import (
	"NotesService/internal/database"
	"NotesService/internal/logger"
	"NotesService/internal/middleware"
	"NotesService/internal/services"
	"github.com/labstack/echo/v4"
)

func main() {
	log := logger.NewLogger(false)

	db, err := database.PostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Успешное подключение к бд")

	svc := services.NewService(db, log)

	router := echo.New()

	api := router.Group("api")

	// api users
	api.POST("/register", svc.Register)
	api.POST("/login", svc.Login)

	// api notes
	api.POST("/note", svc.CreateNoteHandler, middleware.AuthMiddleware)
	api.GET("/notes", svc.GetAllNotesHandler, middleware.AuthMiddleware)
	api.GET("/note/:id", svc.GetNoteHandler, middleware.AuthMiddleware)
	api.PUT("/note/:id", svc.UpdateNoteHandler, middleware.AuthMiddleware)
	api.DELETE("/note/:id", svc.DeleteNoteHandler, middleware.AuthMiddleware)

	router.Logger.Fatal(router.Start(":8080"))

	//!!!!

}
