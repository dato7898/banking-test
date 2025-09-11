package config

import "os"

type Config struct {
	PostgresDSN        string
	KafkaBrokers       []string
	KafkaPaymentsTopic string
	IntegrationAddr    string
}

func Load() *Config {
	return &Config{
		PostgresDSN:        getenv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		KafkaBrokers:       []string{getenv("KAFKA_BROKER", "localhost:9092")},
		KafkaPaymentsTopic: getenv("KAFKA_PAYMENTS_TOPIC", "payments"),
		IntegrationAddr:    getenv("INTEGRATION_ADDR", "localhost:50051"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
