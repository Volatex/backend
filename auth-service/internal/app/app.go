// Пакет app отвечает за конфигурацию и запуск приложения.
package app

import (
	"fmt"
	"gitverse.ru/volatex/backend/pkg/postgres"
	"syscall"

	//"fmt"
	"gitverse.ru/volatex/backend/config"
	"gitverse.ru/volatex/backend/pkg/httpserver"
	"gitverse.ru/volatex/backend/pkg/logger"
	"os"
	"os/signal"
)

// Run инициализирует объекты приложения с помощью конструкторов и запускает его.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Репозиторий: подключение к PostgreSQL.
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// HTTP-сервер: настройка с указанным портом и режимом prefork.
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))

	httpServer.Start()

	// Ожидание сигнала прерывания (Ctrl+C или SIGTERM).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Завершение работы HTTP-сервера.
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
