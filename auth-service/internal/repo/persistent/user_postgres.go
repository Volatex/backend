package persistent

import (
	"context"
	"fmt"

	"gitverse.ru/volatex/backend/internal/entity"
	"gitverse.ru/volatex/backend/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Store(ctx context.Context, user entity.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("email, password").
		Values(user.Email, user.Password).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Pool.QueryRow: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, password").
		From("users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.Builder: %w", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}
