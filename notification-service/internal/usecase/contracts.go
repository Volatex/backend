package usecase

import "context"

type (
	NotificationUseCase interface {
		SendVerificationCode(ctx context.Context, email string) error
	}
)
