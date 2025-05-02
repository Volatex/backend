package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gitverse.ru/volatex/backend/config"
	"gitverse.ru/volatex/backend/internal/controller/http/middleware"
	v1 "gitverse.ru/volatex/backend/internal/controller/http/v1"
	"gitverse.ru/volatex/backend/internal/usecase"
	"gitverse.ru/volatex/backend/pkg/logger"
)

// Swagger spec:
// @title       Authentication and Authorization Service
// @description API for user authentication and authorization
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, u usecase.User, l logger.Interface) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	apiV1Group := app.Group("/v1")
	{
		v1.NewAuthRoutes(apiV1Group, u, l)
	}
}
