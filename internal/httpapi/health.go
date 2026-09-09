package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	resp := HealthResponse{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Ошибка записи ответа:", err)
	}
}
