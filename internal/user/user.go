package user

import (
	"errors"
	"time"
)

var ErrUsernameTaken = errors.New("логин уже занят")

type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}
