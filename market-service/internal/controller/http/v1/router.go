package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gitverse.ru/volatex/backend/market-service/internal/usecase"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

func NewMarketRoutes(apiV1Group fiber.Router, u usecase.Market, l logger.Interface) {
	r := &Market{u: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	strategyGroup := apiV1Group.Group("/strategy")
	{
		strategyGroup.Post("/add", r.saveStrategy)
		strategyGroup.Post("/add_token", r.saveUserToken)
		strategyGroup.Get("/get_strategies", r.getUserStrategies)
	}
}
