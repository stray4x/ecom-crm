package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string `env:"APP_ENV"     validate:"required,oneof=development production"`
	AppDomain        string `env:"APP_DOMAIN"  validate:"required"`
	AppPort          string `env:"APP_PORT"    validate:"required"`
	DBHost           string `env:"DB_HOST"     validate:"required"`
	DBPort           string `env:"DB_PORT"     validate:"required"`
	DBUser           string `env:"DB_USER"     validate:"required"`
	DBPassword       string `env:"DB_PASSWORD" validate:"required"`
	DBName           string `env:"DB_NAME"     validate:"required"`
	JWTAccessSecret  string `env:"JWT_ACCESS_SECRET"  validate:"required"`
	JWTRefreshSecret string `env:"JWT_REFRESH_SECRET" validate:"required"`
	RedisHost        string `env:"REDIS_HOST" validate:"required"`
	RedisPort        string `env:"REDIS_PORT" validate:"required"`
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	cfg := &Config{
		AppEnv:           os.Getenv("APP_ENV"),
		AppDomain:        os.Getenv("APP_DOMAIN"),
		AppPort:          os.Getenv("APP_PORT"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		JWTAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}

	return cfg
}
