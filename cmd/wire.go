//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/kostinp/edu-platform-backend/internal/user/repository"
	http "github.com/kostinp/edu-platform-backend/internal/user/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/kostinp/edu-platform-backend/pkg/config"
	"github.com/kostinp/edu-platform-backend/pkg/db"
	echo "github.com/labstack/echo/v4"
)

// Обёртки для строковых типов
func provideBotToken(cfg *config.Config) config.BotToken {
	return config.BotToken(cfg.Telegram.Token)
}

func provideJwtSecret(cfg *config.Config) config.JwtSecret {
	return config.JwtSecret(cfg.JWT.Secret)
}

func InitializeServer(cfg *config.Config) (*echo.Echo, error) {
	wire.Build(
		// --- DB ---
		db.ConnectPostgres,

		// --- Repositories ---
		repository.NewPostgresUserRepository,
		wire.Bind(new(usecase.UserRepository), new(*repository.PostgresUserRepository)),

		repository.NewPostgresVisitorEventRepo,
		wire.Bind(new(repository.VisitorEventRepository), new(*repository.PostgresVisitorEventRepo)),

		repository.NewPostgresSessionRepository,
		wire.Bind(new(repository.SessionRepository), new(*repository.PostgresSessionRepository)),

		// --- Usecases ---
		usecase.NewSessionUsecase,
		wire.Bind(new(usecase.SessionUsecase), new(*usecase.SessionUsecaseImpl)),

		usecase.NewUserService,
		usecase.NewVisitorEventUsecase,

		// --- Handlers ---
		http.NewUserHandler,
		http.NewVisitorEventHandler,
		http.NewSessionHandler,

		// --- Telegram Auth ---
		provideBotToken,
		provideJwtSecret,
		http.NewTelegramAuthHandler,

		// --- Server ---
		newEchoServer,
	)
	return nil, nil
}

func InitializeSessionUsecase(cfg *config.Config) (usecase.SessionUsecase, error) {
	wire.Build(
		db.ConnectPostgres,
		repository.NewPostgresSessionRepository,
		wire.Bind(new(repository.SessionRepository), new(*repository.PostgresSessionRepository)),
		usecase.NewSessionUsecase,
		wire.Bind(new(usecase.SessionUsecase), new(*usecase.SessionUsecaseImpl)),
	)
	return nil, nil
}
