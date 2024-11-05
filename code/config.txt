package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config хранит параметры конфигурации
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}, nil
}
