package request

type SaveStrategy struct {
	Figi         string  `json:"figi" validate:"required"`
	BuyPrice     float64 `json:"buy_price" validate:"required"`
	BuyQuantity  int     `json:"buy_quantity" validate:"required"`
	SellPrice    float64 `json:"sell_price" validate:"required"`
	SellQuantity int     `json:"sell_quantity" validate:"required"`
}
