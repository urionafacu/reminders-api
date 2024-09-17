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
	ServerPort    string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables")
	}

	config := &Config{
		DBPath:        getEnv("DB_PATH", "reminder.db"),
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailTo:       getEnv("EMAIL_TO", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		APIKey:        getEnv("API_KEY", ""),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
	}

	if config.EmailFrom == "" || config.EmailTo == "" || config.EmailPassword == "" {
		log.Fatal("Missing essential email configuration")
	}

	if config.APIKey == "" {
		log.Fatal("API key is not set")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
