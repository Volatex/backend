package user

import "context"

type (
	NotificationService interface {
		SendVerificationCode(ctx context.Context, email string) error
	}
)
