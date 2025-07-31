//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"

	"github.com/kostinp/edu-platform-backend/internal/shared/db"
	"github.com/kostinp/edu-platform-backend/internal/user/repository"
	http "github.com/kostinp/edu-platform-backend/internal/user/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
)

var UserSet = wire.NewSet(
	// DB
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
	ProvideBotToken,
	ProvideJwtSecret,
	http.NewTelegramAuthHandler,
)
