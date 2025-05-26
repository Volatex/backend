package usecase

import (
	"context"

	"github.com/google/uuid"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
)

type (
	Market interface {
		Create(ctx context.Context, s *entity.Strategy) error
		SaveUserToken(ctx context.Context, token *entity.UserToken) error
		GetUserStrategies(ctx context.Context, userID uuid.UUID) ([]*entity.Strategy, error)
		GetUserStockPositions(ctx context.Context, userID uuid.UUID) ([]*entity.StockPosition, error)
		Delete(ctx context.Context, strategyID int, userID uuid.UUID) error
	}
)
