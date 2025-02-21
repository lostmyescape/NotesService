package logger

import (
	"github.com/labstack/gommon/log"
	"os"
)

func NewLogger(writeToFile bool) *log.Logger {
	// настройка логгера для записи в файл
	logger := log.New("notes")
	if writeToFile {
		logFile, err := os.OpenFile("app.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		logger.SetOutput(logFile)
	}

	logger.SetLevel(log.INFO) // Уровень логирования: DEBUG, INFO, WARN, ERROR
	logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line} ${message}")

	// Пример логирования
	logger.Info("Application started")

	return logger
}
