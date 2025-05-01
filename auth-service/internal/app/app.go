// Пакет app конфигурирует и запускает приложение.
package app

import (
	"gitverse.ru/volatex/backend/config"
	"gitverse.ru/volatex/backend/pkg/logger"
)

// Run создаёт объекты через конструкторы.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
}
