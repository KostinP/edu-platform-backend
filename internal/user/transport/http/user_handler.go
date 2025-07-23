package transport

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/kostinp/edu-platform-backend/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UserService
}

func NewUserHandler(uc usecase.UserService) *UserHandler {
	return &UserHandler{userUsecase: uc}
}

func (h *UserHandler) LinkVisitorToUser(c echo.Context) error {
	visitorIDRaw := c.Get(middleware.VisitorIDKey)
	if visitorIDRaw == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "visitor_id не найден"})
	}

	visitorID, err := uuid.Parse(visitorIDRaw.(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "некорректный visitor_id"})
	}

	userIDParam := c.Param("user_id") // предполагаем, что ты передаёшь user_id в URL
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "некорректный user_id"})
	}

	err = h.userUsecase.LinkVisitorToUser(c.Request().Context(), userID, visitorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "visitor успешно связан с user"})
}

func GetVisitorIDHandler(c echo.Context) error {
	visitorID := c.Get(middleware.VisitorIDKey)
	if visitorID == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "visitor_id not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"visitor_id": visitorID.(string),
		"message":    "Привет, гость!",
	})
}
