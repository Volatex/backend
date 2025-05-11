package email

import (
	"context"
	"fmt"
	"gitverse.ru/volatex/backend/notification-service/internal/repo"
	emailpkg "gitverse.ru/volatex/backend/notification-service/pkg/email"
	"math/rand"
	"time"
)

const codeTTL = 15 * 60

type UseCase struct {
	emailSender emailpkg.SMTPClient
	codeRepo    repo.VerificationCodeRepo
}

func New(emailSender emailpkg.SMTPClient, codeRepo repo.VerificationCodeRepo) *UseCase {
	return &UseCase{
		emailSender: emailSender,
		codeRepo:    codeRepo,
	}
}

func (uc *UseCase) SendVerificationCode(ctx context.Context, email string) error {
	code := generateCode()
	fmt.Println("Generated code:", code)

	if err := uc.codeRepo.Set(ctx, email, code, codeTTL); err != nil {
		fmt.Println("Error saving code to Redis:", err)
		return fmt.Errorf("failed to store code: %w", err)
	}
	fmt.Println("Successfully saved code to Redis")

	subject := "Your verification code"
	body := fmt.Sprintf("Your verification code is: %s", code)

	err := uc.emailSender.Send(email, subject, body)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Successfully sent email to", email)
	return nil
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (uc *UseCase) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	storedCode, err := uc.codeRepo.Get(ctx, email)
	if err != nil {
		if err == repo.ErrCodeNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get code from repo: %w", err)
	}

	return storedCode == code, nil
}
