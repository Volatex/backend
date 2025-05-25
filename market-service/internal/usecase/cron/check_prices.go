package cron

import (
	"context"
	"fmt"

	"gitverse.ru/volatex/backend/market-service/internal/entity"
	"gitverse.ru/volatex/backend/market-service/internal/repo"
	"gitverse.ru/volatex/backend/market-service/pkg/external"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

type CheckPricesUseCase interface {
	Execute(ctx context.Context) error
}

type checkPricesUseCase struct {
	strategyRepo repo.StrategyRepo
	logger       logger.Interface
}

func NewCheckPricesUseCase(strategyRepo repo.StrategyRepo, logger logger.Interface) CheckPricesUseCase {
	return &checkPricesUseCase{
		strategyRepo: strategyRepo,
		logger:       logger,
	}
}

func (uc *checkPricesUseCase) Execute(ctx context.Context) error {
	tokens, err := uc.strategyRepo.GetAllUserTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	strategies, err := uc.strategyRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get strategies: %w", err)
	}

	for _, token := range tokens {
		config := external.TinkoffConfig{
			Token: token.TinkoffToken,
		}

		client, err := external.NewTinkoffClient(config, uc.logger)
		if err != nil {
			uc.logger.Error(err, "Failed to create Tinkoff client", "user_id", token.UserID)
			continue
		}

		var userStrategies []*entity.Strategy
		for _, strategy := range strategies {
			if strategy.UserID == token.UserID {
				userStrategies = append(userStrategies, strategy)
			}
		}

		for _, strategy := range userStrategies {
			price, err := client.GetLastPrice(ctx, strategy.Figi, true)
			if err != nil {
				uc.logger.Error(err, "Failed to get price",
					"user_id", token.UserID,
					"ticker", strategy.Figi)
				continue
			}

			uc.logger.Info("Current price for strategy",
				"user_id", token.UserID,
				"ticker", strategy.Figi,
				"current_price", price,
				"buy_price", strategy.BuyPrice,
				"sell_price", strategy.SellPrice)

			accountID, err := client.GetUserAccount(ctx)
			if err != nil {
				uc.logger.Error(err, "Failed to get user account",
					"user_id", token.UserID)
				continue
			}

			if price <= strategy.BuyPrice {
				uc.logger.Info("Buy condition met, executing buy order",
					"user_id", token.UserID,
					"ticker", strategy.Figi,
					"current_price", price,
					"target_price", strategy.BuyPrice,
					"quantity", strategy.BuyQuantity)

				err = client.BuyAsset(ctx, strategy.Figi, int64(strategy.BuyQuantity), accountID)
				if err != nil {
					uc.logger.Error(err, "Failed to create buy order",
						"user_id", token.UserID,
						"ticker", strategy.Figi)
					continue
				}

				uc.logger.Info("Buy order executed successfully",
					"user_id", token.UserID,
					"ticker", strategy.Figi,
					"quantity", strategy.BuyQuantity,
					"price", price)

				// TODO: Обновить статус стратегии после успешной покупки
				// Например, пометить стратегию как активную или обновить целевую цену продажи
			}

			if price >= strategy.SellPrice {
				uc.logger.Info("Sell condition met, executing sell order",
					"user_id", token.UserID,
					"ticker", strategy.Figi,
					"current_price", price,
					"target_price", strategy.SellPrice,
					"quantity", strategy.SellQuantity)

				err = client.SellAsset(ctx, strategy.Figi, int64(strategy.SellQuantity), accountID)
				if err != nil {
					uc.logger.Error(err, "Failed to create sell order",
						"user_id", token.UserID,
						"ticker", strategy.Figi)
					continue
				}

				uc.logger.Info("Sell order executed successfully",
					"user_id", token.UserID,
					"ticker", strategy.Figi,
					"quantity", strategy.SellQuantity,
					"price", price)

				// TODO: Обновить статус стратегии после успешной продажи
				// Например, пометить стратегию как завершенную или сбросить цены
			}
		}
	}

	return nil
}
