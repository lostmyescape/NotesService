package database

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"os"
)

type DB struct {
	Name     string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
}

func PostgresConnection() (*sql.DB, error) {
	config := getDBConfig()

	// коннект к бд
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)

	// подключение к бд
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Успешное подключение к бд")

	return db, nil
}

// getDBConfig получает данные из config.yaml
func getDBConfig() DB {
	configPath := "../internal/config/config.yaml"

	// читает из конфига
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Ошибка чтения файла yaml: %v", err)
	}

	var config DB

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Ошибка парсинга yaml: %v", err)
	}

	return config

}
