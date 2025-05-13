package entity

import "github.com/google/uuid"

type Strategy struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Figi         string
	BuyPrice     float64
	BuyQuantity  int
	SellPrice    float64
	SellQuantity int
	TinkoffToken string
}
