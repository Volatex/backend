package entity

import "time"

type VerificationCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}
