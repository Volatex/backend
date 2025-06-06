package app

import (
	"fmt"
	"gitverse.ru/volatex/backend/internal/clients/notificationgrpc"
	"syscall"

	"gitverse.ru/volatex/backend/internal/controller/http"
	"gitverse.ru/volatex/backend/internal/repo/persistent"
	"gitverse.ru/volatex/backend/internal/usecase/user"
	"gitverse.ru/volatex/backend/pkg/postgres"

	"os"
	"os/signal"

	"gitverse.ru/volatex/backend/config"
	"gitverse.ru/volatex/backend/pkg/httpserver"
	"gitverse.ru/volatex/backend/pkg/logger"
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

	userRepo := persistent.New(pg)

	// Инициализация gRPC клиента
	notifClient, err := notificationgrpc.New(cfg.GRPC.NotificationAddr)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - notificationgrpc.New: %w", err))
	}
	//defer notifClient.Close()

	userUseCase := user.New(userRepo, notifClient, cfg.JWT.Secret)

	// HTTP-сервер: настройка с указанным портом и режимом prefork.
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))

	http.NewRouter(httpServer.App, cfg, userUseCase, l)

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
