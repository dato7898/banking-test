package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	PostgresDSN string
	JWTSecret   string
	RedisAddr   string
	RedisPass   string
	RedisDb     int
	SessionTTL  time.Duration
}

func Load() *Config {
	redisDb, err := strconv.Atoi(getenv("REDIS_DB", "0"))
	if err != nil {
		log.Fatal("Error reading Redis DB:", err)
	}

	sessionTTL, err := time.ParseDuration(getenv("SESSION_TTL", "3600s"))
	if err != nil {
		log.Fatal("Error reading session TTL:", err)
	}

	return &Config{
		Port:        getenv("PORT", "8081"),
		PostgresDSN: getenv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", "super-secret-key"),
		RedisAddr:   getenv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   getenv("REDIS_PASS", ""),
		RedisDb:     redisDb,
		SessionTTL:  sessionTTL,
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
