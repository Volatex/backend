package cron

import (
	"context"
	"encoding/json"
	"fmt"

	"gitverse.ru/volatex/backend/market-service/internal/repo"
)

type printStrategiesUseCase struct {
	strategyRepo repo.StrategyRepo
}

func NewPrintStrategiesUseCase(strategyRepo repo.StrategyRepo) PrintStrategiesUseCase {
	return &printStrategiesUseCase{
		strategyRepo: strategyRepo,
	}
}

func (uc *printStrategiesUseCase) Execute(ctx context.Context) error {
	strategies, err := uc.strategyRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get strategies: %w", err)
	}

	jsonData, err := json.MarshalIndent(strategies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal strategies: %w", err)
	}

	fmt.Printf("Current user strategies (%d):\n%s\n", len(strategies), string(jsonData))
	return nil
}
