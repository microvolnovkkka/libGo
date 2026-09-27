package user

import "context"

// какие операции с хранилищем нужны сервису.
type Repository interface {
	Create(ctx context.Context, username, passwordHash string) (User, error)
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
