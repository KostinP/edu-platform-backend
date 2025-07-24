package transport

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	SessionUsecase usecase.SessionUsecase
}

func NewSessionHandler(sessionUC usecase.SessionUsecase) *SessionHandler {
	return &SessionHandler{SessionUsecase: sessionUC}
}

// GET /me/sessions
func (h *SessionHandler) ListSessions(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	sessions, err := h.SessionUsecase.ListSessions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get sessions"})
	}

	return c.JSON(http.StatusOK, sessions)
}

// DELETE /me/sessions/:id
func (h *SessionHandler) DeleteSession(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session id"})
	}

	err = h.SessionUsecase.DeleteSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete session"})
	}

	return c.NoContent(http.StatusNoContent)
}

// POST /me/inactivity-timeout
func (h *SessionHandler) SetInactivityTimeout(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	type Request struct {
		TimeoutSeconds int64 `json:"timeout_seconds"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	timeout := time.Duration(req.TimeoutSeconds) * time.Second
	if timeout < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "timeout must be positive"})
	}

	err = h.SessionUsecase.SetInactivityTimeout(c.Request().Context(), userID, timeout)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to set inactivity timeout"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET /me/inactivity-timeout
func (h *SessionHandler) GetInactivityTimeout(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	timeout, err := h.SessionUsecase.GetInactivityTimeout(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get inactivity timeout"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"timeout_seconds": int64(timeout.Seconds()),
	})
}
