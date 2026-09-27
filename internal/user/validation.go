package user

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidUsername = errors.New("некорректный логин")
	ErrInvalidPassword = errors.New("некорректный пароль")
)

// пример " Fedor_12 " → "fedor_12"
func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// проверяет уже нормализованный логин
func ValidateUsername(username string) error {
	// Для разрешённых ASCII-символов число байт равно числу символов
	if len(username) < 3 || len(username) > 32 {
		return ErrInvalidUsername
	}

	for _, ch := range username {
		isLetter := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		isUnderscore := ch == '_'

		if !isLetter && !isDigit && !isUnderscore {
			return ErrInvalidUsername
		}
	}

	return nil
}

// 8–128 рун
// пробелы и кириллица разрешены, результ не меняем
func ValidatePassword(password string) error {
	if !utf8.ValidString(password) {
		return ErrInvalidPassword
	}

	length := utf8.RuneCountInString(password)
	if length < 8 || length > 128 {
		return ErrInvalidPassword
	}

	return nil
}
