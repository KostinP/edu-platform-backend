package main

import (
	"fmt"

	"github.com/kostinp/edu-platform-backend/pkg/config"
	"github.com/kostinp/edu-platform-backend/pkg/db"
	"github.com/kostinp/edu-platform-backend/pkg/logger"
	"github.com/kostinp/edu-platform-backend/pkg/middleware"
	"github.com/kostinp/edu-platform-backend/pkg/telegram"

	"github.com/kostinp/edu-platform-backend/internal/user/http/transport"
	"github.com/kostinp/edu-platform-backend/internal/user/repository"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.Load()

	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	dbPool := db.ConnectPostgres(cfg)

	userRepo := repository.NewPostgresUserRepository(dbPool)
	userService := usecase.NewUserService(userRepo)
	userHandler := transport.NewUserHandler(userService)
	e.POST("/api/users/:user_id/link-visitor", userHandler.LinkVisitorToUser)

	visitorEventRepo := repository.NewPostgresVisitorEventRepo(dbPool)
	visitorEventUsecase := usecase.NewVisitorEventUsecase(visitorEventRepo)
	visitorEventHandler := transport.NewVisitorEventHandler(visitorEventUsecase)
	e.GET("/api/visitor", transport.GetVisitorIDHandler)
	e.POST("/api/visitor/events", visitorEventHandler.LogEvent)

	telegramAuthHandler := transport.NewTelegramAuthHandler(userService, cfg.Telegram.Token)
	e.POST("/api/telegram/auth", telegramAuthHandler.Auth)

	// middleware
	e.Use(middleware.VisitorMiddleware)
	e.Use(middleware.SetUserIDMiddleware)
	e.Use(middleware.LinkVisitorWithUser(userService))

	bot, err := telegram.New(cfg.Telegram.Token)
	if err != nil {
		logger.Fatal("Не удалось запустить Telegram-бота", err)
	}
	_ = bot // пока не используется

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("🚀 Запуск сервера на %s", addr))

	if err := e.Start(addr); err != nil {
		logger.Fatal("❌ Не удалось запустить сервер", err)
	}
}
