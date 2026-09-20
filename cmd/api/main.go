package main

import (
	"errors"
	"libgoproj/internal/config"
	"libgoproj/internal/httpapi"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, nil) //обработчик логов в текстовый вид в Stdout их оформляет
	logger := slog.New(logHandler)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpapi.HealthHandler)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Ошибка конфигурации", "error", err)
		os.Exit(1)
	}
	srv := http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

	logger.Info("Запуск HTTP-сервера", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Ошибка HTTP-сервера", "error", err)

		//пока что ненулевая ошибка
		// - : завершает немедленно без деферов
		os.Exit(1)
	}
}
