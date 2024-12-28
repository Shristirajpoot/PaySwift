package config

import (
	"os"
)

type Config struct {
	DatabaseURL       string
	PaystackSecretKey string
	CourierAPIKey     string
	ArkeselAPIKey     string
	ServerPort        string
	PG_HOST           string
	PG_PORT           string
	PG_NAME           string
	PG_USER           string
	PG_PASS           string
	PG_SSLMODE        string
	DATABASE_URL      string
}

func Load() (*Config, error) {
	return &Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		PaystackSecretKey: os.Getenv("PAYSTACK_SECRET_KEY"),
		CourierAPIKey:     os.Getenv("COURIER_API_KEY"),
		ArkeselAPIKey:     os.Getenv("ARKESEL_API_KEY"),
		ServerPort:        os.Getenv("SERVER_PORT"),
		PG_HOST:           os.Getenv("DB_HOST"),
		PG_USER:           os.Getenv("DB_USER"),
		PG_PASS:           os.Getenv("DB_PASSWORD"),
		PG_NAME:           os.Getenv("DB_NAME"),
		PG_PORT:           os.Getenv("DB_PORT"),
	}, nil
}
