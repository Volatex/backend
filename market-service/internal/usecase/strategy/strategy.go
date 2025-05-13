package strategy

import (
	"context"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
	"gitverse.ru/volatex/backend/market-service/internal/repo"
)

type UseCase struct {
	repo repo.StrategyRepo
}

func New(repo repo.StrategyRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, s *entity.Strategy) error {
	return uc.repo.Store(ctx, s)
}
