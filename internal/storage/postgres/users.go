package postgres

import (
	"context"
	"errors"
	"fmt"

	"libgoproj/internal/user"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(
	ctx context.Context,
	username string,
	passwordHash string) (user.User, error) {
	const query = `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, created_at
	`

	var created user.User

	// Выполняем запрос и записываем возвращённые поля в created
	err := r.pool.QueryRow(ctx, query, username, passwordHash).Scan(
		&created.ID,
		&created.Username,
		&created.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && //если внутри err есть ошибка типа pgErr то запишет в неё
			pgErr.Code == "23505" && //нарушение уникальности
			pgErr.ConstraintName == "users_username_key" {
			return user.User{}, user.ErrUsernameTaken
		}

		return user.User{}, fmt.Errorf("создание пользователя: %w", err)
	}

	return created, nil
}
