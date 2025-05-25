package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitverse.ru/volatex/backend/market-service/config"
	"gitverse.ru/volatex/backend/market-service/internal/controller/http"
	"gitverse.ru/volatex/backend/market-service/internal/repo/persistent"
	"gitverse.ru/volatex/backend/market-service/internal/usecase/cron"
	"gitverse.ru/volatex/backend/market-service/internal/usecase/strategy"
	"gitverse.ru/volatex/backend/market-service/pkg/external"
	"gitverse.ru/volatex/backend/market-service/pkg/httpserver"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
	"gitverse.ru/volatex/backend/market-service/pkg/postgres"
)

// Run инициализирует зависимости и запускает HTTP-сервер.
func Run(cfg *config.Config) {
	// Логгер
	l := logger.New(cfg.Log.Level)

	// Подключение к PostgreSQL
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Репозиторий
	strategyRepo := persistent.NewStrategyRepo(pg)

	// Math service client
	mathClient, err := external.NewMathServiceClient(external.MathServiceConfig{
		Endpoint: cfg.Math.Endpoint,
	}, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - mathService.New: %w", err))
	}

	// Usecase
	marketUseCase := strategy.New(strategyRepo, l, nil, mathClient)

	// Инициализация и запуск CRON job для печати стратегий
	printStrategiesUseCase := cron.NewPrintStrategiesUseCase(strategyRepo)
	printStrategiesJob := cron.NewPrintStrategiesJob(printStrategiesUseCase)
	printStrategiesJob.Start()
	defer printStrategiesJob.Stop()

	// Инициализация и запуск CRON job для проверки цен
	checkPricesUseCase := cron.NewCheckPricesUseCase(strategyRepo, l)
	checkPricesJob := cron.NewCheckPricesJob(checkPricesUseCase)
	checkPricesJob.Start()
	defer checkPricesJob.Stop()

	// HTTP-сервер
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, marketUseCase, l)

	// Запуск HTTP-сервера
	httpServer.Start()

	// Ожидание завершения по Ctrl+C или SIGTERM
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Завершение сервера
	if err := httpServer.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
