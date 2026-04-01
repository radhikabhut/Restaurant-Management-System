package config

import (
	"github.com/caarlos0/env/v10"
)

var RuntimeConfig AppConfig

type AppConfig struct {
	Server   Server
	Database Database
}

type Server struct {
	Port      int    `env:"SERVER_PORT" envDefault:"8080"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"supersecret"`
}

type Database struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
	Name     string `env:"DB_NAME" envDefault:"restaurant"`
}

// LoadConfig populates RuntimeConfig from environment variables.
func LoadConfig() error {
	return env.Parse(&RuntimeConfig)
}
