package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppKey       string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	RedisHost    string
	RedisPort    string
	RedisPass    string
	JwtSecret    string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		// Not fatal — env vars may be set directly in the environment
		fmt.Println("[WARN][CONFIG] No .env file found, falling back to environment variables")
	}

	cfg := &Config{
		AppKey:       getEnv("APP_KEY", ""),
		DatabaseHost: getEnv("DATABASE_HOST", "127.0.0.1"),
		DatabasePort: getEnv("DATABASE_PORT", "3306"),
		DatabaseUser: getEnv("DATABASE_USER", "root"),
		DatabasePass: getEnv("DATABASE_PASSWORD", ""),
		DatabaseName: getEnv("DATABASE_DATABASE", "featherpanel"),
		RedisHost:    getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		RedisPass:    getEnv("REDIS_PASSWORD", ""),
		JwtSecret:    getEnv("JWT_SECRET", "secret"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
