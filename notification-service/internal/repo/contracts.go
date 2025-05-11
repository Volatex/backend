package repo

import "context"

type (
	VerificationCodeRepo interface {
		Set(ctx context.Context, email string, code string, ttlSeconds int) error
		Get(ctx context.Context, email string) (string, error)
		Delete(ctx context.Context, email string) error
	}
)
