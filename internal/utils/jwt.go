package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("supersecretkey")

// GenerateToken генерация jwt токена на 24 часа
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
