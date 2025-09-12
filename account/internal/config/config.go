package config

import "os"

type Config struct {
	PostgresDSN string
	GRPCPort    string
}

func Load() *Config {
	return &Config{
		PostgresDSN: getenv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		GRPCPort:    getenv("GRPC_PORT", "50052"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
