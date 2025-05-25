package entity

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type UserTariff string

const (
	TariffInvestor UserTariff = "investor" // 0.3% commission
	TariffTrader   UserTariff = "trader"   // 0.04% commission
	TariffPremium  UserTariff = "premium"  // 0% commission
)

// Volatility represents annualized volatility as a percentage
type Volatility float64

// MarshalJSON implements custom JSON marshaling for Volatility
func (v Volatility) MarshalJSON() ([]byte, error) {
	// Convert to percentage and round to 2 decimal places
	percentage := float64(v) * 100
	formatted := fmt.Sprintf("%.2f%%", percentage)
	return json.Marshal(formatted)
}

type StockPosition struct {
	UserID       uuid.UUID  `json:"user_id"`
	Figi         string     `json:"figi"`
	Ticker       string     `json:"ticker"`
	Name         string     `json:"name"`
	Quantity     int        `json:"quantity"`
	CurrentPrice float64    `json:"current_price"`
	TotalValue   float64    `json:"total_value"`
	Commission   float64    `json:"commission"`
	Tariff       UserTariff `json:"tariff"`
	Volatility   Volatility `json:"volatility"`
}

// CalculateCommission calculates the commission for selling the entire position
func (sp *StockPosition) CalculateCommission() {
	totalValue := float64(sp.Quantity) * sp.CurrentPrice

	switch sp.Tariff {
	case TariffInvestor:
		sp.Commission = totalValue * 0.003 // 0.3%
	case TariffTrader:
		sp.Commission = totalValue * 0.0004 // 0.04%
	case TariffPremium:
		sp.Commission = 0 // 0%
	default:
		sp.Commission = totalValue * 0.003 // Default to investor tariff
	}

	sp.TotalValue = totalValue
}
