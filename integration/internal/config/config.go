package config

import "os"

type Config struct {
	GRPCPort    string
	ExternalURL string
}

func Load() *Config {
	return &Config{
		GRPCPort:    getenv("GRPC_PORT", "50051"),
		ExternalURL: getenv("EXTERNAL_URL", "http://localhost:9000/stub/payments"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
