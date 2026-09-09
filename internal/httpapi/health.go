package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	resp := healthResponse{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		//потом будем logger по указателю передавать с мейна
		log.Println("Ошибка записи ответа:", err)
	}
}
