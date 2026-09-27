package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      uint32 = 19 * 1024 // КБ: примерно 19 МБ памяти
	argonIterations  uint32 = 2         // Количество проходов алгоритма
	argonParallelism uint8  = 1         // Степень параллелизма
	argonSaltLength         = 16        // Длина случайной соли в байтах
	argonHashLength  uint32 = 32        // Длина результата в байтах
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
