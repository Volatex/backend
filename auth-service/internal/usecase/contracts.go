package usecase

import (
	"context"

	"gitverse.ru/volatex/backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	User interface {
		Register(ctx context.Context, user entity.User) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
	}
)
