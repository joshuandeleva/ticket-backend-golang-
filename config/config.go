package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
}

func NewEnvConfig() *EnvConfig {
	err := godotenv.Load()
	if err!= nil {
		log.Fatal("Error loading.env file %e" , err)
	}
	config := &EnvConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatal("Error parsing.env file %e" , err)
	}
	return config
}