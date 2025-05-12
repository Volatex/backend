package response

import "github.com/google/uuid"

type User struct {
	Id    uuid.UUID `json:"id" example:"1"`
	Email string    `json:"email" example:"user@example.com"`
}
