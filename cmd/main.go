package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/kostinp/edu-platform-backend/docs"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/kostinp/edu-platform-backend/pkg/config"
	"github.com/kostinp/edu-platform-backend/pkg/logger"
)

// @title Edu Platform API
// @version 1.0
// @description Backend for Edu Platform with gamification
// @termsOfService https://edu-platform.com/terms
// @contact.name Support Team
// @contact.email support@edu-platform.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	server, err := InitializeServer(cfg)
	if err != nil {
		logger.Fatal("Ошибка инициализации сервера", err)
	}

	// Инициализируем usecase отдельно
	sessionUsecase, err := InitializeSessionUsecase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем задачу очистки
	go usecase.StartSessionCleanupTask(context.Background(), sessionUsecase, time.Hour*24)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("🚀 Запуск сервера на %s", addr))

	if err := server.Start(addr); err != nil {
		logger.Fatal("❌ Не удалось запустить сервер", err)
	}
}
