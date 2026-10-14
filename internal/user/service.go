package user

import (
    "context"
    "errors"
    "fmt"
)

// какие операции с хранилищем нужны сервису.
type Repository interface {
	Create(ctx context.Context, username, passwordHash string) (User, error)
	FindByUsername(ctx context.Context, username string) (Credentials, error)
}

type Service struct {
	repo Repository
}

// Получаем готовый репозиторий и сохраняем его внутри сервиса.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(
	ctx context.Context,
	username string,
	password string,
) (User, error) {
	username = NormalizeUsername(username)

	if err := ValidateUsername(username); err != nil {
		return User{}, err
	}

	// HashPassword сначала проверяет пароль, затем создаёт хеш.
	passwordHash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	// В хранилище передаём уже подготовленные данные.
	return s.repo.Create(ctx, username, passwordHash)
}

// Authenticate проверяет логин и пароль.
// При успехе возвращает пользователя без хеша пароля.
func (s *Service) Authenticate(
    ctx context.Context,
    username string,
    password string,
) (User, error) {
    username = NormalizeUsername(username)

    credentials, err := s.repo.FindByUsername(ctx, username)
    if errors.Is(err, ErrUserNotFound) {
        return User{}, ErrInvalidCredentials
    }
    if err != nil {
        return User{}, fmt.Errorf(
            "получение пользователя для входа: %w",
            err,
        )
    }

    match, err := VerifyPassword(password, credentials.PasswordHash)
    if err != nil {
        return User{}, fmt.Errorf(
            "проверка сохранённого хеша пароля: %w",
            err,
        )
    }
    if !match {
        return User{}, ErrInvalidCredentials
    }

    return credentials.User, nil
}