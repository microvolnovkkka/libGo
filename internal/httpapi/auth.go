package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"time"

	"libgoproj/internal/user"
)

type AuthHandler struct {
	service *user.Service
	logger  *slog.Logger
}

func NewAuthHandler(service *user.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// Формат входящего JSON.
type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Формат успешного ответа.
type registerResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil || mediaType != "application/json" {
		h.writeJSON(w, http.StatusUnsupportedMediaType, errorResponse{Error: "нужен Content-Type: application/json"})
		return
	}

	// Для логина и пароля достаточно тела запроса до 4 КиБ.
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.writeJSON(w, http.StatusRequestEntityTooLarge, errorResponse{Error: "тело запроса слишком большое"})
		} else {
			h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "не удалось прочитать тело запроса"})
		}
		return
	}

	var req registerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, errorResponse{Error: "некорректный JSON"})
		return
	}

	// Все правила регистрации выполняет сервис.
	created, err := h.service.Register(
		r.Context(),
		req.Username,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidUsername):
			h.writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "логин: от 3 до 32 символов, только a-z, 0-9 и _",
			})

		case errors.Is(err, user.ErrInvalidPassword):
			h.writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: "пароль должен содержать от 8 до 128 символов",
			})

		case errors.Is(err, user.ErrUsernameTaken):
			h.writeJSON(w, http.StatusConflict, errorResponse{
				Error: "логин уже занят",
			})

		default:
			h.logger.Error("Ошибка регистрации пользователя", "error", err)

			h.writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "внутренняя ошибка сервера",
			})
		}
		return
	}

	h.writeJSON(w, http.StatusCreated, registerResponse{
		ID:        created.ID,
		Username:  created.Username,
		CreatedAt: created.CreatedAt,
	})
}

// Общая отправка JSON для методов этого обработчика.
func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		h.logger.Error("Ошибка формирования JSON-ответа", "error", err)

		status = http.StatusInternalServerError
		data = []byte(`{"error":"внутренняя ошибка сервера"}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		h.logger.Error("Ошибка отправки HTTP-ответа", "error", err)
	}
}
