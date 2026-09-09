package config

import "os"

type Config struct {
	HTTPAddr string
}

func Load() Config {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "127.0.0.1:8080"
	}
	return Config{
		HTTPAddr: httpAddr,
	}
}
