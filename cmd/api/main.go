package main

import (
	"context"
	"errors"
	"fmt"
	"libgoproj/internal/config"
	"libgoproj/internal/httpapi"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("загрузка конфигурации: %w", err)
	}
	// отмена ctx не закрывает пул соединений, ctx нужен чтобы Ping занимал не более 5 сек
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("создание пула PostgreSQL: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("проверка подключения к PostgreSQL: %w", err)
	}
	cancel()
	logger.Info("Подключение к PostgreSQL установлено")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpapi.HealthHandler)

	srv := http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

	logger.Info("Запуск HTTP-сервера", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("работа HTTP-сервера: %w", err)

	}
	return nil
}

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, nil) //обработчик логов в текстовый вид в Stdout их оформляет
	logger := slog.New(logHandler)
	if err := run(logger); err != nil {
		logger.Error("Ошибка приложения", "error", err)
		os.Exit(1)
	}
}
