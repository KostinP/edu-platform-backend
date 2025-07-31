package transport

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/middleware"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

type VisitorEventHandler struct {
	usecase usecase.VisitorEventUsecase
}

func NewVisitorEventHandler(uc usecase.VisitorEventUsecase) *VisitorEventHandler {
	return &VisitorEventHandler{usecase: uc}
}

type LogEventRequest struct {
	EventType string         `json:"event_type" validate:"required"`
	EventData map[string]any `json:"event_data"`
}

// @Summary Логировать событие посетителя
// @Tags Visitors
// @Accept json
// @Produce json
// @Param event body LogEventRequest true "Событие"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /visitor/events [post]
func (h *VisitorEventHandler) LogEvent(c echo.Context) error {
	visitorIDRaw := c.Get(middleware.VisitorIDKey)
	if visitorIDRaw == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "visitor_id отсутствует"})
	}

	visitorID, err := uuid.Parse(visitorIDRaw.(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid visitor_id"})
	}

	req := new(LogEventRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат запроса"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = h.usecase.LogEvent(c.Request().Context(), visitorID, req.EventType, req.EventData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось сохранить событие"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "событие сохранено"})
}
