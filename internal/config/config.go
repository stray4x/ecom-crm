package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string `env:"APP_PORT"    validate:"required"`
	DBHost     string `env:"DB_HOST"     validate:"required"`
	DBPort     string `env:"DB_PORT"     validate:"required"`
	DBUser     string `env:"DB_USER"     validate:"required"`
	DBPassword string `env:"DB_PASSWORD" validate:"required"`
	DBName     string `env:"DB_NAME"     validate:"required"`
}

func initConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	cfg := &Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}

	return cfg
}

var Env = initConfig()
