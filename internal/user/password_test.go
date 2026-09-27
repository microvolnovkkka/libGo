package user

import (
	"errors"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	first, err := HashPassword("study-go-password")
	if err != nil {
		t.Fatalf("первое хеширование: %v", err)
	}

	second, err := HashPassword("study-go-password")
	if err != nil {
		t.Fatalf("второе хеширование: %v", err)
	}

	if !strings.HasPrefix(first, "$argon2id$v=19$") {
		t.Error("результат должен содержать алгоритм Argon2id и его версию")
	}

	if first == second {
		t.Error("повторное хеширование должно использовать новую соль")
	}
}

func TestHashPasswordRejectsInvalidPassword(t *testing.T) {
	hash, err := HashPassword("short")

	if !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("получили ошибку %v, ожидали ErrInvalidPassword", err)
	}

	if hash != "" {
		t.Error("для недопустимого пароля не должен возвращаться хеш")
	}
}
