package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	GRPCPort    string
	HTTPPort    string
	DatabaseUrl string
	JWTSecret   string
	NatsURL     string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		AppEnv:      getEnv("APP_ENV", "local"),
		GRPCPort:    getEnv("GRPC_PORT", "8081"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseUrl: getEnv("DB_DSN", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		NatsURL:     getEnv("NATS_URL", "nats://localhost:4222"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
