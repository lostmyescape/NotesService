package main

import (
	"NotesService/internal/database"
	"NotesService/internal/logger"
	"github.com/labstack/echo/v4"
)

func main() {
	log := logger.NewLogger(false)

	_, err := database.PostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	router := echo.New()

	_ = router.Group("api")

}
