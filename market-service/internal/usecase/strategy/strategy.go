package strategy

import (
	"context"
	"github.com/google/uuid"
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

func (uc *UseCase) SaveUserToken(ctx context.Context, token *entity.UserToken) error {
	return uc.repo.StoreUserToken(ctx, token)
}

func (uc *UseCase) GetUserStrategies(ctx context.Context, userID uuid.UUID) ([]*entity.Strategy, error) {
	return uc.repo.GetByUserID(ctx, userID)
}
