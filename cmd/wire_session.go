//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/kostinp/edu-platform-backend/internal/shared/db"
	"github.com/kostinp/edu-platform-backend/internal/user/repository"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
)

var SessionUsecaseSet = wire.NewSet(
	db.ConnectPostgres,
	repository.NewPostgresSessionRepository,
	wire.Bind(new(repository.SessionRepository), new(*repository.PostgresSessionRepository)),
	usecase.NewSessionUsecase,
	wire.Bind(new(usecase.SessionUsecase), new(*usecase.SessionUsecaseImpl)),
)
