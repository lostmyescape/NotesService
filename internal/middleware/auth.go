package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

var jwtSecret = []byte("your-secret-key")

// AuthMiddleware аутентификация по jwt токену
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// берем передаваемый токен и проверяем пустой ли он
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
		}

		// удаляем префикс
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
		}

		// парсим токен и проверяем подпись
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token sign"})
		}

		// извлекаем userID из claims и добавляем в контекст
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["userID"] == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid claims"})
		}

		// сохраняем userID в контексте
		c.Set("userID", int(claims["userID"].(float64)))

		return next(c)
	}
}
