package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const VisitorIDKey = "visitor_id"

func VisitorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(VisitorIDKey)
		if err != nil || cookie.Value == "" {
			// Создаем новый visitor_id
			newVisitorID := uuid.NewString()
			cookie = &http.Cookie{
				Name:     VisitorIDKey,
				Value:    newVisitorID,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   365 * 24 * 60 * 60, // 1 год
				Secure:   false,              // поменять на true при HTTPS
				SameSite: http.SameSiteLaxMode,
			}
			c.SetCookie(cookie)
		}

		// Добавляем visitor_id в контекст Echo
		c.Set(VisitorIDKey, cookie.Value)

		return next(c)
	}
}
