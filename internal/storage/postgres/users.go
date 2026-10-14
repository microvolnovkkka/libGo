package postgres

import (
	"context"
	"errors"
	"fmt"

	"libgoproj/internal/user"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
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

// FindByUsername ожидает логин, уже нормализованный сервисом.
func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (user.Credentials, error) {
	const query = `
        SELECT id, username, created_at, password_hash
        FROM users
        WHERE username = $1
    `

	var credentials user.Credentials

	err := r.pool.QueryRow(ctx, query, username).Scan(
		&credentials.User.ID,
		&credentials.User.Username,
		&credentials.User.CreatedAt,
		&credentials.PasswordHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Credentials{}, user.ErrUserNotFound
	}
	if err != nil {
		return user.Credentials{}, fmt.Errorf(
			"поиск пользователя по логину: %w",
			err,
		)
	}

	return credentials, nil
}
