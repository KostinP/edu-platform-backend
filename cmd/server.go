package main

import (
	"github.com/kostinp/edu-platform-backend/pkg/config"
	"github.com/kostinp/edu-platform-backend/pkg/middleware"

	transport "github.com/kostinp/edu-platform-backend/internal/user/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"

	echo "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewServer создает и конфигурирует Echo-сервер
func newEchoServer(
	cfg *config.Config,
	userHandler *transport.UserHandler,
	visitorEventHandler *transport.VisitorEventHandler,
	telegramAuthHandler *transport.TelegramAuthHandler,
	sessionHandler *transport.SessionHandler,
	sessionUsecase usecase.SessionUsecase,
	userService *usecase.UserService,
) (*echo.Echo, error) {
	e := echo.New()

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Middleware
	jwtMiddleware := middleware.JWTMiddleware([]byte(cfg.JWT.Secret), sessionUsecase)
	e.Use(jwtMiddleware)
	e.Use(middleware.VisitorMiddleware)
	e.Use(middleware.SetUserIDMiddleware)
	e.Use(middleware.LinkVisitorWithUser(userService))

	// Роуты
	e.POST("/api/users/:user_id/link-visitor", userHandler.LinkVisitorToUser)
	e.GET("/api/visitor", transport.GetVisitorIDHandler)
	e.POST("/api/visitor/events", visitorEventHandler.LogEvent)
	e.POST("/api/telegram/auth", telegramAuthHandler.Auth)

	// Защищённые роуты /me/* для сессий
	meGroup := e.Group("/me", jwtMiddleware)
	meGroup.GET("/sessions", sessionHandler.ListSessions)
	meGroup.DELETE("/sessions/:id", sessionHandler.DeleteSession)
	meGroup.POST("/inactivity-timeout", sessionHandler.SetInactivityTimeout)
	meGroup.GET("/inactivity-timeout", sessionHandler.GetInactivityTimeout)

	return e, nil
}
