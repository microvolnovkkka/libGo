package main

import (
	"errors"
	"libgoproj/internal/httpapi"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpapi.HealthHandler)
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Ошибка сервера:", err)
	}
}
