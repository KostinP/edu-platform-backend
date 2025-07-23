package transport

import (
	"net/http"

	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/kostinp/edu-platform-backend/pkg/telegram"
	"github.com/labstack/echo/v4"
)

type TelegramAuthHandler struct {
	userService *usecase.UserService
	botToken    string
}

func NewTelegramAuthHandler(userService *usecase.UserService, botToken string) *TelegramAuthHandler {
	return &TelegramAuthHandler{userService: userService, botToken: botToken}
}

func (h *TelegramAuthHandler) Auth(c echo.Context) error {
	authData := telegram.ParseTelegramAuth(c.Request())

	if !telegram.VerifyTelegramAuth(authData, h.botToken) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid telegram auth"})
	}

	user, err := h.userService.GetByTelegramID(c.Request().Context(), authData.ID)
	if err != nil {
		// Создание нового юзера, если не найден
		user, err = h.userService.CreateFromTelegramAuth(c.Request().Context(), authData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать пользователя"})
		}
	}

	// Возвращаем user_id для сессионного хранения (в куки, в localstorage и т.д.)
	return c.JSON(http.StatusOK, map[string]any{
		"user_id":   user.ID.String(),
		"username":  user.Username,
		"full_name": user.FullName,
		"role":      user.Role,
	})
}
