package repo

import (
	"context"

	"gitverse.ru/volatex/backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	UserRepo interface {
		Store(ctx context.Context, user entity.User) error
		GetByEmail(ctx context.Context, email string) (entity.User, error)
	}
)
