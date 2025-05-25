package entity

import (
	"time"

	"github.com/google/uuid"
)

type Strategy struct {
	ID           int
	UserID       uuid.UUID
	Figi         string
	BuyPrice     float64
	BuyQuantity  int
	SellPrice    float64
	SellQuantity int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
