package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath        string
	EmailFrom     string
	EmailTo       string
	EmailPassword string
	APIKey        string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables")
	}

	return &Config{
		DBPath:        getEnv("DB_PATH", "reminder.db"),
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailTo:       getEnv("EMAIL_TO", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		APIKey:        getEnv("API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
