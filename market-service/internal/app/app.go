package app

import (
	"fmt"
	"gitverse.ru/volatex/backend/market-service/config"
	"gitverse.ru/volatex/backend/market-service/internal/controller/http"
	"gitverse.ru/volatex/backend/market-service/internal/repo/persistent"
	"gitverse.ru/volatex/backend/market-service/internal/usecase/strategy"
	"gitverse.ru/volatex/backend/market-service/pkg/httpserver"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
	"gitverse.ru/volatex/backend/market-service/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
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

	// Репозиторий и usecase
	strategyRepo := persistent.NewStrategyRepo(pg)
	marketUseCase := strategy.New(strategyRepo)

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
