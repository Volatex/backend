package v1

import (
	"github.com/go-playground/validator/v10"
	"gitverse.ru/volatex/backend/market-service/internal/usecase"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

type Market struct {
	u usecase.Market
	l logger.Interface
	v *validator.Validate
}
