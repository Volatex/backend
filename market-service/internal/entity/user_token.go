package entity

import "github.com/google/uuid"

type UserToken struct {
	UserID       uuid.UUID `json:"user_id"`
	TinkoffToken string    `json:"tinkoff_token"`
}
