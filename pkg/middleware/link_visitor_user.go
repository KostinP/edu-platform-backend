package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

func LinkVisitorWithUser(userUC *usecase.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			visitorIDRaw := c.Get(VisitorIDKey)
			userIDRaw := c.Get(UserIDKey) // Предполагаем, что у тебя уже есть UserID в контексте

			// если одного из них нет — ничего не делаем
			if visitorIDRaw == nil || userIDRaw == nil {
				return next(c)
			}

			visitorID, err := uuid.Parse(visitorIDRaw.(string))
			if err != nil {
				return next(c) // некорректный visitorID — просто пропускаем
			}
			userID, ok := userIDRaw.(uuid.UUID)
			if !ok {
				return next(c)
			}

			// пытаемся связать — в случае ошибки логгируем, но не прерываем запрос
			err = userUC.LinkVisitorToUser(context.Background(), userID, visitorID)
			if err != nil {
				// логировать стоит в будущем
			}

			return next(c)
		}
	}
}
