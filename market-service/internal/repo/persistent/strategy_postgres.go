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
	sql, args, err := r.Builder.
		Insert("user_strategies").
		Columns("user_id", "figi", "buy_price", "buy_quantity", "sell_price", "sell_quantity").
		Values(s.UserID, s.Figi, s.BuyPrice, s.BuyQuantity, s.SellPrice, s.SellQuantity).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("StrategyRepo - Store - Build: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)
	if err := row.Scan(&s.ID); err != nil {
		return fmt.Errorf("StrategyRepo - Store - Scan: %w", err)
	}

	return nil
}

func (r *StrategyRepo) StoreUserToken(ctx context.Context, t *entity.UserToken) error {
	sql, args, err := r.Builder.
		Insert("user_tokens").
		Columns("user_id", "tinkoff_token").
		Values(t.UserID, t.TinkoffToken).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET tinkoff_token = EXCLUDED.tinkoff_token, updated_at = NOW()").
		ToSql()
	if err != nil {
		return fmt.Errorf("StrategyRepo - StoreUserToken - Build: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("StrategyRepo - StoreUserToken - Exec: %w", err)
	}

	return nil
}

func (r *StrategyRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Strategy, error) {
	sql, args, err := r.Builder.
		Select("id", "user_id", "figi", "buy_price", "buy_quantity", "sell_price", "sell_quantity", "created_at", "updated_at").
		From("user_strategies").
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StrategyRepo - GetByUserID - Build: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("StrategyRepo - GetByUserID - Query: %w", err)
	}
	defer rows.Close()

	var strategies []*entity.Strategy
	for rows.Next() {
		var s entity.Strategy
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Figi, &s.BuyPrice, &s.BuyQuantity,
			&s.SellPrice, &s.SellQuantity, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("StrategyRepo - GetByUserID - Scan: %w", err)
		}
		strategies = append(strategies, &s)
	}
	return strategies, nil
}

func (r *StrategyRepo) GetAll(ctx context.Context) ([]*entity.Strategy, error) {
	sql, args, err := r.Builder.
		Select("id", "user_id", "figi", "buy_price", "buy_quantity", "sell_price", "sell_quantity", "created_at", "updated_at").
		From("user_strategies").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StrategyRepo - GetAll - Build: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("StrategyRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()

	var strategies []*entity.Strategy
	for rows.Next() {
		var s entity.Strategy
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Figi, &s.BuyPrice, &s.BuyQuantity,
			&s.SellPrice, &s.SellQuantity, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("StrategyRepo - GetAll - Scan: %w", err)
		}
		strategies = append(strategies, &s)
	}
	return strategies, nil
}

func (r *StrategyRepo) GetAllUserTokens(ctx context.Context) ([]*entity.UserToken, error) {
	sql, args, err := r.Builder.
		Select("user_id", "tinkoff_token").
		From("user_tokens").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StrategyRepo - GetAllUserTokens - Build: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("StrategyRepo - GetAllUserTokens - Query: %w", err)
	}
	defer rows.Close()

	var tokens []*entity.UserToken
	for rows.Next() {
		var t entity.UserToken
		err := rows.Scan(
			&t.UserID,
			&t.TinkoffToken,
		)
		if err != nil {
			return nil, fmt.Errorf("StrategyRepo - GetAllUserTokens - Scan: %w", err)
		}
		tokens = append(tokens, &t)
	}
	return tokens, nil
}
