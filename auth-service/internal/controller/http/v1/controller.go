package v1

import (
	"github.com/go-playground/validator/v10"
	"gitverse.ru/volatex/backend/internal/usecase"
	"gitverse.ru/volatex/backend/pkg/logger"
)

type Auth struct {
	u usecase.User
	l logger.Interface
	v *validator.Validate
}
