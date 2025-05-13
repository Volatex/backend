package persistent

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
	"gitverse.ru/volatex/backend/market-service/pkg/postgres"
)

type StrategyRepo struct {
	*postgres.Postgres
}

func NewStrategyRepo(pg *postgres.Postgres) *StrategyRepo {
	return &StrategyRepo{pg}
}

func (r *StrategyRepo) Store(ctx context.Context, s *entity.Strategy) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	sql, args, err := r.Builder.
		Insert("user_strategies").
		Columns("id", "user_id", "figi", "buy_price", "buy_quantity", "sell_price", "sell_quantity", "tinkoff_token").
		Values(s.ID, s.UserID, s.Figi, s.BuyPrice, s.BuyQuantity, s.SellPrice, s.SellQuantity, s.TinkoffToken).
		ToSql()
	if err != nil {
		return fmt.Errorf("StrategyRepo - Store - Build: %w", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("StrategyRepo - Store - Exec: %w", err)
	}
	return nil
}
