package config

import "os"

type Config struct {
	PostgresDSN        string
	KafkaBrokers       []string
	KafkaPaymentsTopic string
	HTTPPort           string
	JWTSecret          string
	AccountAddr        string
}

func Load() *Config {
	return &Config{
		PostgresDSN:        getenv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		KafkaBrokers:       []string{getenv("KAFKA_BROKER", "localhost:9092")},
		KafkaPaymentsTopic: getenv("KAFKA_PAYMENTS_TOPIC", "payments"),
		HTTPPort:           getenv("HTTP_PORT", "8080"),
		JWTSecret:          getenv("JWT_SECRET", "super-secret-key"),
		AccountAddr:        getenv("ACCOUNT_ADDR", "localhost:50052"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
