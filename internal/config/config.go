package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (Config, error) {
	httpAddr := os.Getenv("HTTP_ADDR")
	if strings.TrimSpace(httpAddr) == "" {
		httpAddr = "127.0.0.1:8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databaseURL) == "" {
		return Config{}, errors.New("не задана переменная окружения DATABASE_URL")
	}
	return Config{
		HTTPAddr:    httpAddr,
		DatabaseURL: databaseURL,
	}, nil
}
