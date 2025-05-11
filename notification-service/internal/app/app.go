// Пакет app отвечает за конфигурацию и запуск notification-сервиса.
package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"gitverse.ru/volatex/backend/notification-service/config"
	"gitverse.ru/volatex/backend/notification-service/internal/controller/grpc"
	v1 "gitverse.ru/volatex/backend/notification-service/internal/controller/grpc/v1"
	"gitverse.ru/volatex/backend/notification-service/internal/repo/redis"
	usecaseEmail "gitverse.ru/volatex/backend/notification-service/internal/usecase/email"
	pkgEmail "gitverse.ru/volatex/backend/notification-service/pkg/email"
	"gitverse.ru/volatex/backend/notification-service/pkg/logger"
	redisPkg "gitverse.ru/volatex/backend/notification-service/pkg/redis"

	"google.golang.org/grpc/reflection"
)

// Run инициализирует объекты приложения с помощью конструкторов и запускает его.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Подключение к Redis
	rdb, err := redisPkg.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer rdb.Close()

	// Инициализация email сендера
	emailSender := pkgEmail.NewSMTPClient(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	// Репозиторий и usecase
	codeRepo := redis.NewCodeRepo(rdb.Client)
	useCase := usecaseEmail.New(*emailSender, codeRepo)

	// gRPC хендлер и сервер
	handler := v1.NewNotificationHandler(useCase)
	server := grpc.NewServer(l, handler)

	reflection.Register(server)

	// Создание TCP-листенера
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - net.Listen: %w", err))
	}

	// Запуск сервера в отдельной горутине
	go func() {
		l.Info("Starting gRPC server on :" + cfg.GRPC.Port)
		if err := server.Serve(lis); err != nil {
			l.Fatal(fmt.Errorf("app - Run - server.Serve: %w", err))
		}
	}()

	// Ожидание сигнала завершения
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	}

	// Завершение работы сервера
	server.GracefulStop()
	l.Info("gRPC server stopped gracefully")
}
