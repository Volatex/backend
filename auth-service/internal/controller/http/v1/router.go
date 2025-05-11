package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gitverse.ru/volatex/backend/internal/usecase"
	"gitverse.ru/volatex/backend/pkg/logger"
)

func NewAuthRoutes(apiV1Group fiber.Router, u usecase.User, l logger.Interface) {
	r := &Auth{u: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	authGroup := apiV1Group.Group("/auth")

	{
		authGroup.Post("/register", r.register)
		authGroup.Post("/verify-email", r.verifyEmail)
	}
}
