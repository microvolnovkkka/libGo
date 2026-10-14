package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      uint32 = 19 * 1024 // КБ: примерно 19 МБ памяти
	argonIterations  uint32 = 2         // Количество проходов алгоритма
	argonParallelism uint8  = 1         // Степень параллелизма
	argonSaltLength         = 16        // Длина случайной соли в байтах
	argonHashLength  uint32 = 32        // Длина результата в байтах
)

var ErrInvalidPasswordHash = errors.New(
	"некорректный или неподдерживаемый хеш пароля",
)

func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	salt := make([]byte, argonSaltLength)
	rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		argonHashLength,
	)

	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonIterations,
		argonParallelism,
		saltBase64,
		hashBase64,
	)

	return encoded, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	// Наш формат укладывается в этот размер.
	if len(encodedHash) > 256 {
		return false, ErrInvalidPasswordHash
	}

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return false, ErrInvalidPasswordHash
	}

	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		return false, ErrInvalidPasswordHash
	}

	// Читаем параметры, использованные при создании хеша.
	var memory, iterations uint32
	var parallelism uint8

	n, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	)
	if err != nil || n != 3 {
		return false, ErrInvalidPasswordHash
	}

	// Проверяем всю строку параметров, включая отсутствие лишнего текста.
	expectedParams := fmt.Sprintf(
		"m=%d,t=%d,p=%d",
		memory,
		iterations,
		parallelism,
	)
	if parts[3] != expectedParams {
		return false, ErrInvalidPasswordHash
	}

	// Допустимые затраты для проверки пароля в нашем приложении.
	if parallelism < 1 || parallelism > 4 ||
		iterations < 1 || iterations > 3 ||
		memory < 8*uint32(parallelism) || memory > 64*1024 {
		return false, ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil || len(salt) != argonSaltLength {
		return false, ErrInvalidPasswordHash
	}

	storedHash, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil || len(storedHash) != int(argonHashLength) {
		return false, ErrInvalidPasswordHash
	}

	// Повторяем вычисление с введённым паролем и сохранёнными параметрами.
	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(storedHash)),
	)

	return subtle.ConstantTimeCompare(actualHash, storedHash) == 1, nil
}
