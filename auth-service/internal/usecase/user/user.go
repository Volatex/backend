package user

import (
	"context"
	"fmt"

	"gitverse.ru/volatex/backend/internal/entity"
	"gitverse.ru/volatex/backend/internal/repo"
)

type UseCase struct {
	repo     repo.UserRepo
	notifier NotificationService
}

func New(r repo.UserRepo, notifier NotificationService) *UseCase {
	return &UseCase{
		repo:     r,
		notifier: notifier,
	}
}

func (uc *UseCase) Register(ctx context.Context, user entity.User) (entity.User, error) {
	err := uc.repo.Store(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	if err := uc.notifier.SendVerificationCode(ctx, user.Email); err != nil {
		fmt.Printf("UserUseCase - Register - uc.notifier.SendVerificationCode: %v\n", err)
	}

	return user, nil
}

func (uc *UseCase) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - GetByEmail - uc.repo.GetByEmail: %w", err)
	}

	return user, nil
}

func (uc *UseCase) VerifyEmail(ctx context.Context, email string, code string) error {
	// TODO: Сделать валидацию кода

	err := uc.repo.SetEmailVerified(ctx, email)
	if err != nil {
		return fmt.Errorf("UserUseCase - VerifyEmail - uc.repo.MarkEmailVerified: %w", err)
	}
	return nil
}
