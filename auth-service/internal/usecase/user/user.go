package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"

	"gitverse.ru/volatex/backend/internal/entity"
	"gitverse.ru/volatex/backend/internal/repo"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrIncorrectPassword  = errors.New("password incorrect")
)

type UseCase struct {
	repo      repo.UserRepo
	notifier  NotificationService
	jwtSecret string
}

func New(r repo.UserRepo, notifier NotificationService, jwtSecret string) *UseCase {
	return &UseCase{
		repo:      r,
		notifier:  notifier,
		jwtSecret: jwtSecret,
	}
}

func (uc *UseCase) Register(ctx context.Context, user entity.User) (entity.User, error) {
	err := uc.repo.Store(ctx, &user)
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
	valid, err := uc.notifier.VerifyCode(ctx, email, code)
	if err != nil {
		return fmt.Errorf("UserUseCase - VerifyEmail - notifier.VerifyCode: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid verification code")
	}

	err = uc.repo.SetEmailVerified(ctx, email)
	if err != nil {
		return fmt.Errorf("UserUseCase - VerifyEmail - uc.repo.SetEmailVerified: %w", err)
	}
	return nil
}

func (uc *UseCase) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("SignIn - GetByEmail: %w", err)
	}

	if !verifyPassword(user.Password, password) {
		return "", ErrIncorrectPassword
	}

	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("SignIn - token.SignedString: %w", err)
	}

	return tokenString, nil
}

func verifyPassword(hashedPassword, password string) bool {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false
	}

	memory := uint32(65536)
	iterations := uint32(1)
	parallelism := uint8(4)

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, 32)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	return hashB64 == parts[5]
}
