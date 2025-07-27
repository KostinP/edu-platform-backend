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

	// Swagger - публичный, без JWT
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health - публичный, без JWT
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Ваш JWT middleware
	jwtMiddleware := middleware.JWTMiddleware([]byte(cfg.JWT.Secret), sessionUsecase)

	// Другие middleware, которые нужны для всех запросов
	e.Use(middleware.VisitorMiddleware)
	e.Use(middleware.SetUserIDMiddleware)
	e.Use(middleware.LinkVisitorWithUser(userService))

	// Публичные маршруты без JWT
	e.POST("/api/users/:user_id/link-visitor", userHandler.LinkVisitorToUser)
	e.GET("/api/visitor", transport.GetVisitorIDHandler)
	e.POST("/api/visitor/events", visitorEventHandler.LogEvent)
	e.POST("/api/telegram/auth", telegramAuthHandler.Auth)

	// Создаем группу для маршрутов, защищённых JWT
	apiProtected := e.Group("/api")
	apiProtected.Use(jwtMiddleware)

	// Добавьте сюда все защищённые роуты, например:
	apiProtected.GET("/me/sessions", sessionHandler.ListSessions)
	apiProtected.DELETE("/me/sessions/:id", sessionHandler.DeleteSession)
	apiProtected.POST("/me/inactivity-timeout", sessionHandler.SetInactivityTimeout)
	apiProtected.GET("/me/inactivity-timeout", sessionHandler.GetInactivityTimeout)

	return e, nil
}
