package user

import "github.com/kostinp/edu-platform-backend/internal/shared/config"

func ProvideBotToken(cfg *config.Config) config.BotToken {
	return config.BotToken(cfg.Telegram.Token)
}

func ProvideJwtSecret(cfg *config.Config) config.JwtSecret {
	return config.JwtSecret(cfg.JWT.Secret)
}
