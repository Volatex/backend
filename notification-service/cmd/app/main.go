package main

import (
	"gitverse.ru/volatex/backend/notification-service/config"
	"gitverse.ru/volatex/backend/notification-service/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config errors: %s", err)
	}

	app.Run(cfg)
}
