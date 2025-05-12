package entity

import "github.com/google/uuid"

type User struct {
	Id              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	IsEmailVerified bool      `json:"is_email_verified"`
}
