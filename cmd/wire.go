//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/user"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	echo "github.com/labstack/echo/v4"
)

// Главный wire-компонент
func InitializeServer(cfg *config.Config) (*echo.Echo, error) {
	wire.Build(
		user.UserSet,
		newEchoServer,
	)
	return nil, nil
}

// Для middleware и background job'ов
func InitializeSessionUsecase(cfg *config.Config) (usecase.SessionUsecase, error) {
	wire.Build(SessionUsecaseSet)
	return nil, nil
}
