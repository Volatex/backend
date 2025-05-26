package repo

import (
	"context"

	"github.com/google/uuid"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
)

type (
	StrategyRepo interface {
		Store(ctx context.Context, s *entity.Strategy) error
		StoreUserToken(ctx context.Context, t *entity.UserToken) error
		GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Strategy, error)
		GetByID(ctx context.Context, id int) (*entity.Strategy, error)
		GetAll(ctx context.Context) ([]*entity.Strategy, error)
		GetAllUserTokens(ctx context.Context) ([]*entity.UserToken, error)
		Delete(ctx context.Context, id int) error
	}
)
