package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "gitverse.ru/volatex/backend/market-service/docs"

	"gitverse.ru/volatex/backend/market-service/config"
	"gitverse.ru/volatex/backend/market-service/internal/controller/http/middleware"
	v1 "gitverse.ru/volatex/backend/market-service/internal/controller/http/v1"
	"gitverse.ru/volatex/backend/market-service/internal/usecase"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Market Service
// @description API for managing trading strategies
// @version     1.0
// @host        localhost:8081
// @BasePath    /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(app *fiber.App, cfg *config.Config, market usecase.Market, l logger.Interface) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	app.Static("/swagger", "./docs")

	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	apiV1Group := app.Group("/v1")
	{
		v1.NewMarketRoutes(apiV1Group, market, l)
	}
}
