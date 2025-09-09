package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ApiPort    string `env:"API_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_Host"`
	DBPort     int    `env:"DB_Port"`
	DBUser     string `env:"DB_User"`
	DBName     string `env:"DB_Name"`
	DBPassword string `env:"DB_Password"`
	LogLevel   string `env:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	if err := godotenv.Load("./.env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
