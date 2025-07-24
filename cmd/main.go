package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/user/usecase" // именно как пакет
	"github.com/kostinp/edu-platform-backend/pkg/config"
	"github.com/kostinp/edu-platform-backend/pkg/logger"
)

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
