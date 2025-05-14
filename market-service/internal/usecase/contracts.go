package usecase

import (
	"context"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
)

type (
	Market interface {
		Create(ctx context.Context, s *entity.Strategy) error
		SaveUserToken(ctx context.Context, token *entity.UserToken) error
	}
)
