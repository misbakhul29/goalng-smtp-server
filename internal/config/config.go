package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string
	ServerPort   string
	APIKey       string
}

func Load() (*Config, error) {
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}

	smtpPort, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT value: %w", err)
	}

	cfg := &Config{
		SMTPHost:     requireEnv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPEmail:    requireEnv("SMTP_EMAIL"),
		SMTPPassword: requireEnv("SMTP_PASSWORD"),
		ServerPort:   getEnvOrDefault("SERVER_PORT", "8080"),
		APIKey:       requireEnv("API_KEY"),
	}

	return cfg, nil
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return val
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
