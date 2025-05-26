package strategy

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
	"gitverse.ru/volatex/backend/market-service/internal/repo"
	"gitverse.ru/volatex/backend/market-service/pkg/external"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

type UseCase struct {
	repo          repo.StrategyRepo
	logger        logger.Interface
	tinkoffClient external.TinkoffClient
	mathClient    external.MathServiceClient
}

func New(repo repo.StrategyRepo, logger logger.Interface, tinkoffClient external.TinkoffClient, mathClient external.MathServiceClient) *UseCase {
	return &UseCase{
		repo:          repo,
		logger:        logger,
		tinkoffClient: tinkoffClient,
		mathClient:    mathClient,
	}
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

func (uc *UseCase) GetUserStockPositions(ctx context.Context, userID uuid.UUID) ([]*entity.StockPosition, error) {
	tokens, err := uc.repo.GetAllUserTokens(ctx)
	if err != nil {
		uc.logger.Error(err, "Failed to get user tokens")
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	uc.logger.Info("Retrieved tokens", "count", len(tokens))

	var userToken *entity.UserToken
	for _, token := range tokens {
		if token.UserID == userID {
			userToken = token
			break
		}
	}
	if userToken == nil {
		uc.logger.Error(nil, "No token found for user", "user_id", userID)
		return nil, fmt.Errorf("no token found for user")
	}

	uc.logger.Info("Found user token", "user_id", userID, "token_length", len(userToken.TinkoffToken))

	config := external.TinkoffConfig{
		Token: userToken.TinkoffToken,
	}
	client, err := external.NewTinkoffClient(config, uc.logger)
	if err != nil {
		uc.logger.Error(err, "Failed to create Tinkoff client", "user_id", userID)
		return nil, fmt.Errorf("failed to create Tinkoff client: %w", err)
	}

	uc.logger.Info("Tinkoff client created successfully", "user_id", userID)

	// Get user's account ID
	accountID, err := client.GetUserAccount(ctx)
	if err != nil {
		uc.logger.Error(err, "Failed to get user account", "user_id", userID)
		return nil, fmt.Errorf("failed to get user account: %w", err)
	}

	uc.logger.Info("Got user account", "user_id", userID, "account_id", accountID)

	// Get user's portfolio
	portfolio, err := client.GetPortfolio(ctx, accountID)
	if err != nil {
		uc.logger.Error(err, "Failed to get portfolio", "user_id", userID, "account_id", accountID)
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	uc.logger.Info("Got user portfolio", "user_id", userID, "positions_count", len(portfolio.Positions))

	// Get user's tariff (this should be implemented in your system)
	// For now, we'll use a default tariff
	userTariff := entity.TariffInvestor

	positions := make([]*entity.StockPosition, 0, len(portfolio.Positions))
	for _, pos := range portfolio.Positions {
		uc.logger.Info("Processing position", "user_id", userID, "figi", pos.Figi)

		// Get current price for the instrument
		price, err := client.GetLastPrice(ctx, pos.Figi, false)
		if err != nil {
			uc.logger.Warn("Skipping position due to missing price", "user_id", userID, "figi", pos.Figi, "error", err)
			continue
		}

		uc.logger.Info("Got price", "user_id", userID, "figi", pos.Figi, "price", price)

		// Get instrument details
		instrument, err := client.GetInstrumentByFigi(ctx, pos.Figi)
		if err != nil {
			uc.logger.Error(err, "Failed to get instrument details", "user_id", userID, "figi", pos.Figi)
			return nil, fmt.Errorf("failed to get instrument details for %s: %w", pos.Figi, err)
		}

		uc.logger.Info("Got instrument details", "user_id", userID, "figi", pos.Figi, "ticker", instrument.Ticker)

		// Convert Quotation to int
		quantity := int(pos.Quantity.Units)
		if pos.Quantity.Nano > 0 {
			// Round up if there are any nano units
			quantity++
		}

		position := &entity.StockPosition{
			UserID:       userID,
			Figi:         pos.Figi,
			Ticker:       instrument.Ticker,
			Name:         instrument.Name,
			Quantity:     quantity,
			CurrentPrice: price,
			Tariff:       userTariff,
		}

		// Calculate volatility using historical prices
		// Get prices for the last 12 months
		to := time.Now()
		from := to.AddDate(-1, 0, 0) // 12 months ago
		prices, err := client.GetHistoricalPrices(ctx, pos.Figi, from, to)
		if err != nil {
			uc.logger.Warn("Failed to get historical prices for volatility calculation",
				"user_id", userID,
				"figi", pos.Figi,
				"error", err)
			// Continue without volatility
			position.CalculateCommission()
			positions = append(positions, position)
			continue
		}

		// Calculate monthly returns
		if len(prices) < 2 {
			uc.logger.Warn("Not enough historical prices for volatility calculation",
				"user_id", userID,
				"figi", pos.Figi,
				"prices_count", len(prices))
			position.CalculateCommission()
			positions = append(positions, position)
			continue
		}

		// Calculate returns
		returns := make([]float64, len(prices)-1)
		for i := 1; i < len(prices); i++ {
			if prices[i-1] == 0 {
				uc.logger.Warn("Found zero price in historical data",
					"user_id", userID,
					"figi", pos.Figi,
					"index", i-1)
				continue
			}
			returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
		}

		// Log some statistics about returns
		if len(returns) > 0 {
			var sum float64
			var count int
			for _, r := range returns {
				if !math.IsNaN(r) && !math.IsInf(r, 0) {
					sum += r
					count++
				}
			}
			if count > 0 {
				uc.logger.Info("Returns statistics",
					"user_id", userID,
					"figi", pos.Figi,
					"total_returns", len(returns),
					"valid_returns", count,
					"avg_return", sum/float64(count))
			}
		}

		// Calculate volatility using math service
		volatility, err := uc.mathClient.CalculateVolatility(ctx, returns)
		if err != nil {
			uc.logger.Warn("Failed to calculate volatility",
				"user_id", userID,
				"figi", pos.Figi,
				"error", err,
				"returns_count", len(returns))
			// Set volatility to 0 instead of NaN
			position.Volatility = entity.Volatility(0)
			position.CalculateCommission()
			positions = append(positions, position)
			continue
		}

		position.Volatility = entity.Volatility(volatility)
		position.CalculateCommission()
		positions = append(positions, position)

		uc.logger.Info("Position processed",
			"user_id", userID,
			"figi", pos.Figi,
			"quantity", quantity,
			"price", price,
			"volatility", volatility)
	}

	return positions, nil
}

// Delete deletes a strategy by its ID, but only if it belongs to the specified user
func (uc *UseCase) Delete(ctx context.Context, strategyID int, userID uuid.UUID) error {
	// First get the strategy to verify ownership
	strategy, err := uc.repo.GetByID(ctx, strategyID)
	if err != nil {
		uc.logger.Error(err, "Failed to get strategy for deletion",
			"strategy_id", strategyID,
			"user_id", userID)
		return fmt.Errorf("failed to get strategy: %w", err)
	}

	// Verify that the strategy belongs to the user
	if strategy.UserID != userID {
		uc.logger.Warn("Attempt to delete strategy of another user",
			"strategy_id", strategyID,
			"user_id", userID,
			"strategy_user_id", strategy.UserID)
		return fmt.Errorf("strategy not found")
	}

	// Delete the strategy
	if err := uc.repo.Delete(ctx, strategyID); err != nil {
		uc.logger.Error(err, "Failed to delete strategy",
			"strategy_id", strategyID,
			"user_id", userID)
		return fmt.Errorf("failed to delete strategy: %w", err)
	}

	uc.logger.Info("Strategy deleted successfully",
		"strategy_id", strategyID,
		"user_id", userID)
	return nil
}
