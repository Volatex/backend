package strategy

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitverse.ru/volatex/backend/market-service/internal/entity"
	"gitverse.ru/volatex/backend/market-service/internal/repo"
	"gitverse.ru/volatex/backend/market-service/pkg/external"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

type UseCase struct {
	repo   repo.StrategyRepo
	logger logger.Interface
}

func New(repo repo.StrategyRepo, logger logger.Interface) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
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

		position.CalculateCommission()
		positions = append(positions, position)

		uc.logger.Info("Position processed", "user_id", userID, "figi", pos.Figi, "quantity", quantity, "price", price)
	}

	return positions, nil
}
