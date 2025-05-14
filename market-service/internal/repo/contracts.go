package repo

import (
	"context"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
)

type (
	StrategyRepo interface {
		Store(ctx context.Context, s *entity.Strategy) error
		StoreUserToken(ctx context.Context, t *entity.UserToken) error
	}
)
