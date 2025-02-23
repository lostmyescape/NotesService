package main

import (
	"NotesService/internal/database"
	"NotesService/internal/logger"
	"NotesService/internal/services"
	"github.com/labstack/echo/v4"
)

func main() {
	log := logger.NewLogger(false)

	db, err := database.PostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	svc := services.NewService(db, log)

	router := echo.New()

	api := router.Group("api")

	// api users
	api.POST("/register", svc.Register)
	api.POST("/login", svc.Login)
	//api.DELETE("/words/:id", svc.DeleteUser)

	router.Logger.Fatal(router.Start(":8080"))

}
