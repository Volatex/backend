package entity

import (
	"github.com/google/uuid"
	"time"
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
