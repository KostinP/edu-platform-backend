// pkg/di/wire.go
package di

import (
	"github.com/google/wire"
	"github.com/kostinp/pkg/config"
	"github.com/kostinp/pkg/db"
	"github.com/kostinp/pkg/logger"
	"github.com/kostinp/pkg/telemetry"
)

var ProviderSet = wire.NewSet(
	config.Load,
	logger.New,
	db.NewPostgres,
	db.NewRedis,
	telemetry.NewRedisTracker,
	// Добавляем провайдеры для всех сервисов (user, course и т.д.)
)
