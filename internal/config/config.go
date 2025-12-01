package config

import (
	"os"

	"github.com/joho/godotenv"

)

type Config struct {
	DBUrl          string
	RedisAddr      string
	Port           string
	JWTSecret      string
	AzureConnStr   string
	AzureContainer string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	return &Config{
		DBUrl:          os.Getenv("DATABASE_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		Port:           os.Getenv("PORT"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AzureConnStr:   os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
		AzureContainer: os.Getenv("AZURE_CONTAINER_NAME"),
	}
}