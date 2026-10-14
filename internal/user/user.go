package user

import (
	"errors"
	"time"
)

var (
	ErrUsernameTaken      = errors.New("логин уже занят")
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidCredentials = errors.New("неверный логин или пароль")
)

type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}

// Данные, которые нужны сервису для проверки входа.
type Credentials struct {
	User         User
	PasswordHash string `json:"-"`
}
