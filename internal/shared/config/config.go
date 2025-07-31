package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	DB       DBConfig     `yaml:"db"`
	Telegram Telegram     `yaml:"telegram"`
	JWT      JWTConfig    `yaml:"jwt"`
	Mode     string
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

func Load() *Config {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}

	path := "./configs/" + mode + ".yaml"
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Не удалось прочитать конфигурационный файл: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("Не удалось распарсить конфигурационный файл: %v", err)
	}

	cfg.Mode = mode
	return &cfg
}
