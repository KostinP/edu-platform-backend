package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/pkg/config"
	"github.com/kostinp/edu-platform-backend/pkg/logger"
)

func ConnectPostgres(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Fatal("Не удалось подключиться к базе данных", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		logger.Fatal("База данных недоступна", err)
	}

	logger.Info("Подключение к базе данных установлено")
	return dbPool
}
